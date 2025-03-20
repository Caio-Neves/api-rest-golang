package main

import (
	"fmt"
	"net/http"
	"rest-api-example/config"
	"rest-api-example/handlers"
	"rest-api-example/repositories"
	"rest-api-example/routes"
	"rest-api-example/service"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	cfg, err := config.ReadConfigFile("./config/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/ecommerce.log", cfg.Logs),
		MaxSize:    10,
		MaxBackups: 10,
	})
	log.Info("Setup log file successfully")

	dbInstance, err := config.NewDatabaseConnectionPostgreSQL(cfg.PostgresServerDatabase)
	if err != nil {
		log.Fatal(err)
	}
	defer dbInstance.Close()
	log.Info("Database connection established")

	categoryRepository := repositories.NewCategoryRepositoryPostgres(dbInstance)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepository := repositories.NewProductRepositoryPostgres(dbInstance)
	productService := service.NewProductService(productRepository, categoryRepository)
	productHandler := handlers.NewProductHandler(productService)
	log.Info("Repositories and services initialized")

	r := mux.NewRouter()
	routes.InitCategoryRoutes(r, categoryHandler)
	routes.InitProductRoutes(r, productHandler)
	routes.InitAdminProductsRoutes(r, productHandler)
	routes.InitAdminCategoriesRoutes(r, categoryHandler)
	log.Info("Routes initialized")

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      r,
		ErrorLog:     nil,
	}
	log.Info("Server started on port 8080")
	server.ListenAndServe()
}
