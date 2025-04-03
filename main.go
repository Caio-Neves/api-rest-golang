package main

import (
	"fmt"
	"net/http"
	"rest-api-example/category"
	"rest-api-example/config"
	"rest-api-example/product"
	"rest-api-example/routes"
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
	log.Info("Server started on port 8080")
	server.ListenAndServe()
}
