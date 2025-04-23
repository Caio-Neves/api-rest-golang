package product

import (
	"net/http"
	"rest-api-example/auth"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupProductsRoutes(mux *mux.Router, h ProductHandler, authService auth.AuthService) {
	admin := mux.PathPrefix("/admin/products").Subrouter()
	admin.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	admin.Use(authService.AuthenticationMiddleware)
	admin.HandleFunc("",
		middlewares.ValidateSupportedMediaTypes(([]string{"application/json"}),
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.CreateProduct))).Methods(http.MethodOptions,
		http.MethodPost)
	admin.HandleFunc("/_delete",
		middlewares.ValidateSupportedMediaTypes([]string{"application/json"}, h.DeleteProducts)).Methods(http.MethodOptions,
		http.MethodPost)
	admin.HandleFunc("/{id}",
		middlewares.ValidateSupportedMediaTypes(([]string{"application/json"}),
			middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.UpdateProductsFields))).Methods(http.MethodOptions,
		http.MethodPatch)
	admin.HandleFunc("/{id}", h.DeleteProductById).Methods(http.MethodOptions, http.MethodDelete)

	r := mux.PathPrefix("/products").Subrouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.HandleFunc("", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetAllProducts)).Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc("/{id}", middlewares.ValidadeAcceptHeader([]string{"application/json"},
		h.GetProductById)).Methods(http.MethodOptions, http.MethodGet)
}
