package handlers

import (
	"net/http"
	"rest-api-example/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: s,
	}
}

func (h ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: GetAllProducts"))
}

func (h ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: GetProductById"))
}

func (h ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: CreateProduct"))
}

func (h ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: DeleteProducts"))
}

func (h ProductHandler) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: UpdateProducts"))
}

func (h ProductHandler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: DeleteProductById"))
}
