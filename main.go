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

	log.Println(cfg.PostgresServerDatabase.Database)
	log.Println(cfg.PostgresServerDatabase.Host)
	log.Println(cfg.PostgresServerDatabase.User)
	log.Println(cfg.PostgresServerDatabase.Pass)
	log.Println(cfg.PostgresServerDatabase.Port)
	// dbInstance, err := config.NewDatabaseConnectionSqlServer(cfg.SqlServerDatabase)
	dbInstance, err := config.NewDatabaseConnectionPostgreSQL(cfg.PostgresServerDatabase)
	if err != nil {
		log.Fatal(err)
	}

	// categoryRepository := repositories.NewCategoryRepositorySqlServer(dbInstance)
	categoryRepository := repositories.NewCategoryRepositoryPostgres(dbInstance)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepository := repositories.NewProductRepositoryPostgres(dbInstance)
	productService := service.NewProductService(productRepository, categoryRepository)
	productHandler := handlers.NewProductHandler(productService)

	mux := mux.NewRouter()
	routes.InitCategoryRoutes(mux, categoryHandler)
	routes.InitProductRoutes(mux, productHandler)

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      mux,
		ErrorLog:     nil,
	}
	server.ListenAndServe()
}
