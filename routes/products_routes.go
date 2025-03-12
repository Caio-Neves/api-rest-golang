package routes

import (
	"net/http"
	"rest-api-example/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitProductRoutes(mux *mux.Router, h handlers.ProductHandler) {
	r := mux.PathPrefix("/products").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("", h.GetAllProducts).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/{id}", h.GetProductById).Methods(http.MethodOptions, http.MethodGet)
}
