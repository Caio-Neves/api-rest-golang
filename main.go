package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rest-api-example/auth"
	"rest-api-example/category"
	"rest-api-example/config"
	"rest-api-example/product"
	"rest-api-example/user"
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

	r := mux.NewRouter()

	userRepository := user.NewUserRepository(dbInstance)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)
	user.SetupUserRoutes(r, userHandler)

	authService := auth.NewAuthService(userRepository, os.Getenv("SECRET_KEY"))
	authHandler := auth.NewAuthHandler(authService)
	auth.SetupAuthRoutes(r, authHandler, userHandler)

	categoryRepository := category.NewCategoryRepositoryPostgres(dbInstance)
	categoryService := category.NewCategoryService(categoryRepository)
	categoryHandler := category.NewCategoryHandler(categoryService)
	category.SetupCategoriesRoutes(r, categoryHandler, authService)

	productRepository := product.NewProductRepositoryPostgres(dbInstance)
	productService := product.NewProductService(productRepository, categoryRepository)
	productHandler := product.NewProductHandler(productService)
	product.SetupProductsRoutes(r, productHandler, authService)

	log.Info("Successfully initialized all system layers")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      r,
		ErrorLog:     nil,
	}
	log.Info("Server configured")

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
