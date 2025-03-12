package main

import (
	"fmt"
	"log"
	"net/http"
	"rest-api-example/config"
	"rest-api-example/handlers"
	"rest-api-example/repositories"
	"rest-api-example/routes"
	"rest-api-example/service"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	cfg, err := config.ReadConfigFile("./config/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	logPath := fmt.Sprintf("%s/ecommerce.log", cfg.Logs)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 10,
	}
	log.SetOutput(lumberjackLogger)
	log.Println("Setup lumberjack logger")

	dbInstance, err := config.NewDatabaseConnectionPostgreSQL(cfg.PostgresServerDatabase)
	if err != nil {
		log.Fatal(err)
	}

	categoryRepository := repositories.NewCategoryRepositoryPostgres(dbInstance)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepository := repositories.NewProductRepositoryPostgres(dbInstance)
	productService := service.NewProductService(productRepository, categoryRepository)
	productHandler := handlers.NewProductHandler(productService)

	r := mux.NewRouter()
	routes.InitCategoryRoutes(r, categoryHandler)
	routes.InitProductRoutes(r, productHandler)
	routes.InitAdminProductsRoutes(r, productHandler)
	routes.InitAdminCategoriesRoutes(r, categoryHandler)

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      r,
		ErrorLog:     nil,
	}
	server.ListenAndServe()
}
