package routes

import (
	"net/http"
	"rest-api-example/handlers"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
)

func InitProductRoutes(mux *mux.Router, h *handlers.ProductHandler) {
	mux.HandleFunc("/products", h.GetAllProducts).Methods(http.MethodGet)
	mux.HandleFunc("/products/{id}", h.GetProductById).Methods(http.MethodGet)
	mux.HandleFunc("/products", middlewares.ValidateSupportedMediaTypes(([]string{"application/json"}), h.CreateProduct)).Methods(http.MethodPost)
	mux.HandleFunc("/products/_delete", middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.DeleteProducts)).Methods(http.MethodPost)
	mux.HandleFunc("/products/{id}", h.UpdateProducts).Methods(http.MethodPut)
	mux.HandleFunc("/products/{id}", h.DeleteProductById).Methods(http.MethodDelete)
}
