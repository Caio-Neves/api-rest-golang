package routes

import (
	"net/http"
	"rest-api-example/category"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitAdminCategoriesRoutes(m *mux.Router, h category.CategoryHandler) {
	r := m.PathPrefix("/admin/categories").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5500"},
		AllowedMethods: []string{"POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"},
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.CreateCategory))).Methods(http.MethodOptions,
		http.MethodPost)
	r.HandleFunc("/_delete",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.DeleteCategories)).Methods(http.MethodOptions,
		http.MethodPost)
	r.HandleFunc("/{id}",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"},
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.UpdateCategoryFields))).Methods(http.MethodOptions,
		http.MethodPatch)
	r.HandleFunc("/{id}", h.DeleteCategoryById).Methods(http.MethodOptions, http.MethodDelete)
}
