package routes

import (
	"net/http"
	"rest-api-example/product"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitProductRoutes(mux *mux.Router, h product.ProductHandler) {
	r := mux.PathPrefix("/products").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("", h.GetAllProducts).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/{id}", h.GetProductById).Methods(http.MethodOptions, http.MethodGet)
}
