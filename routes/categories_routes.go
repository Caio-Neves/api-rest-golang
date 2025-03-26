package routes

import (
	"net/http"
	"rest-api-example/category"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitCategoryRoutes(mux *mux.Router, h category.CategoryHandler) {
	r := mux.PathPrefix("/categories").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5500"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetAllCategories)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/{id}", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetCategoryById)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/{id}/products", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetAllProductsByCategory)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/_get",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"},
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.GetCategoriesByIds))).Methods(http.MethodOptions,
		http.MethodPost)
}
