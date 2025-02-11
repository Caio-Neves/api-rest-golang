package routes

import (
	"net/http"
	"rest-api-example/handlers"

	"github.com/gorilla/mux"
)

func InitCategoryRoutes(mux *mux.Router, h *handlers.CategoryHandler) {
	mux.HandleFunc("/categories", h.GetAllCategories).Methods(http.MethodGet)
	mux.HandleFunc("/categories/{id}", h.GetCategoryById).Methods(http.MethodGet)
	mux.HandleFunc("/categories/{id}/products", h.GetAllProductsByCategory).Methods(http.MethodGet)
	mux.HandleFunc("/categories", h.CreateCategory).Methods(http.MethodPost)
	mux.HandleFunc("/categories/_delete", h.DeleteCategories).Methods(http.MethodPost)
	mux.HandleFunc("/categories/{id}", h.UpdateCategories).Methods(http.MethodPut)
	mux.HandleFunc("/categories/{id}", h.DeleteCategoryById).Methods(http.MethodDelete)
}
