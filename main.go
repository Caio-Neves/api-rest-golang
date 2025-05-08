package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rest-api-example/auth"
	"rest-api-example/category"
	"rest-api-example/config"
	"rest-api-example/product"
	"rest-api-example/user"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	action := flag.String("action", "run", "defines action to perform: install, uninstall, run, stop")
	configDir := flag.String("configs", os.Getenv("ECOM_CONFIG_DIR"), "path to config directory")
	flag.Parse()

	if configDir == nil || *configDir == "" {
		panic("Config directory is not set")
	}

	cfg, err := config.ReadConfigFile(*configDir)
	if err != nil {
		panic(err)
	}

	if cfg.ServiceSettings.Name == "" || cfg.ServiceSettings.DisplayName == "" || cfg.ServiceSettings.Description == "" {
		panic("Service settings are not set in the config file")
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

	program := &program{
		server:          server,
		serviceSettings: cfg.ServiceSettings,
	}

	var args []string
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-action") {
			args = append(args, strings.Trim(arg, `"`))
		}
	}

	err = program.SetupServerWindows(*action, args)
	if err != nil {
		panic(err)
	}
}

func InitWebApplication(server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server closed under request: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	<-s

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Info("Server shutdown gracefully")
}
