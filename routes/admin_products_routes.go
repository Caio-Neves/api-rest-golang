package routes

import (
	"net/http"
	"rest-api-example/handlers"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitAdminProductsRoutes(mux *mux.Router, h handlers.ProductHandler) {
	r := mux.PathPrefix("/admin/products").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("",
		middlewares.ValidateSupportedMediaTypes(([]string{"application/json"}), h.CreateProduct)).Methods(http.MethodOptions,
		http.MethodPost)
	r.HandleFunc("/_delete",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.DeleteProducts)).Methods(http.MethodOptions,
		http.MethodPost)
	r.HandleFunc("/{id}",
		middlewares.ValidateSupportedMediaTypes(([]string{"application/json"}), h.UpdateProductsFields)).Methods(http.MethodOptions,
		http.MethodPatch)
	r.HandleFunc("/{id}", h.DeleteProductById).Methods(http.MethodOptions, http.MethodDelete)
}
