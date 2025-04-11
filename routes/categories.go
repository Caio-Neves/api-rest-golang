package routes

import (
	"net/http"
	"rest-api-example/category"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupCategoriesRoutes(mux *mux.Router, h category.CategoryHandler) {
	admin := mux.PathPrefix("/admin/categories").Subrouter()
	admin.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	admin.HandleFunc("",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"},
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.CreateCategory))).Methods(http.MethodOptions,
		http.MethodPost)
	admin.HandleFunc("/_delete",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.DeleteCategories)).Methods(http.MethodOptions,
		http.MethodPost)
	admin.HandleFunc("/{id}",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"},
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.UpdateCategoryFields))).Methods(http.MethodOptions,
		http.MethodPatch)
	admin.HandleFunc("/{id}", h.DeleteCategoryById).Methods(http.MethodOptions, http.MethodDelete)

	r := mux.PathPrefix("/categories").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetPaginateCategories)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/{id}", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetCategoryById)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/{id}/products", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetAllProductsByCategory)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/_get",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"},
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.GetCategoriesByIds))).Methods(http.MethodOptions,
		http.MethodPost)
}
