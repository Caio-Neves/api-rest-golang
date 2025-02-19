package routes

import (
	"net/http"
	"rest-api-example/handlers"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
)

func InitCategoryRoutes(mux *mux.Router, h *handlers.CategoryHandler) {
	mux.HandleFunc("/categories", h.GetAllCategories).Methods(http.MethodGet)
	mux.HandleFunc("/categories/{id}", h.GetCategoryById).Methods(http.MethodGet)
	mux.HandleFunc("/categories/_get", middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.GetCategoriesByIds)).Methods(http.MethodPost)
	mux.HandleFunc("/categories/{id}/products", h.GetAllProductsByCategory).Methods(http.MethodGet)
	mux.HandleFunc("/categories", middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.CreateCategory)).Methods(http.MethodPost)
	mux.HandleFunc("/categories/_delete", middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.DeleteCategories)).Methods(http.MethodPost)
	mux.HandleFunc("/categories/{id}", middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.UpdateCategoryFields)).Methods(http.MethodPatch)
	mux.HandleFunc("/categories/{id}", h.DeleteCategoryById).Methods(http.MethodDelete)
}
