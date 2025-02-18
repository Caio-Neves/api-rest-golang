package handlers

import (
	"context"
	"log"
	"net/http"
	errorsApi "rest-api-example/errors"
	"rest-api-example/service"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	categories, err := h.productService.GetAllProducts(ctx, queryParams)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(categories) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: categories,
			Meta: map[string]interface{}{
				"pagination": map[string]int{
					"page":    getQueryInt(queryParams, "page", 1),
					"limit":   getQueryInt(queryParams, "limit", 10),
					"results": len(categories),
				},
			},
		},
	}, http.StatusOK, w)
}

func (h ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		SendJsonError(JsonResponseError{
			Payload: ResponseError{
				Error: errorsApi.ErrUuidInvalido.Error(),
			},
		}, http.StatusBadRequest, w)
		return
	}
	product, err := h.productService.GetProductById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if product.IsEmpty() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: product,
			Meta: map[string]interface{}{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusOK, w)
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
