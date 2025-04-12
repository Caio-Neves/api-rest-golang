package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rest-api-example/category"
	"rest-api-example/config"
	"rest-api-example/product"
	"rest-api-example/routes"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	cfg, err := config.ReadConfigFile("./config/config.toml")
	if err != nil {
		panic(err)
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/ecommerce.log", cfg.Logs),
		MaxSize:    10,
		MaxBackups: 10,
	})
	log.Info("Setup log file successfully")

	dbInstance, err := config.NewDatabaseConnectionPostgreSQL(cfg.PostgresServerDatabase)
	if err != nil {
		panic(err)
	}
	defer dbInstance.Close()
	log.Info("Database connection established")

	categoryRepository := category.NewCategoryRepositoryPostgres(dbInstance)
	productRepository := product.NewProductRepositoryPostgres(dbInstance)
	log.Info("Successful setup of repositories")

	categoryService := category.NewCategoryService(categoryRepository)
	productService := product.NewProductService(productRepository, categoryRepository)
	log.Info("Successful setup of services")

	categoryHandler := category.NewCategoryHandler(categoryService)
	productHandler := product.NewProductHandler(productService)
	log.Info("Successful setup of handlers")

	r := mux.NewRouter()
	routes.SetupProductsRoutes(r, productHandler)
	routes.SetupCategoriesRoutes(r, categoryHandler)
	log.Info("Routes initialized")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      r,
		ErrorLog:     nil,
	}
	log.Info("Server configured")

	go func() {
		if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server closed under request: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	<-s

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Info("Server shutdown gracefully")
}
