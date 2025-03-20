package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/service"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(s service.ProductService) ProductHandler {
	return ProductHandler{
		productService: s,
	}
}

func (h ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	categories, err := h.productService.GetAllProducts(ctx, queryParams)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: categories,
			Meta: map[string]any{
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
	op := "ProductHandler.GetProductById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	product, err := h.productService.GetProductById(ctx, id)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: product,
			Meta: map[string]any{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusOK, w)
}

func (h ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.CreateProduct()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var product entities.Product
	log.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	product, err = h.productService.CreateProduct(ctx, product)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: product,
			Meta: map[string]any{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusCreated, w)
}

func (h ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.DeleteProducts()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var idsString []string
	err := json.NewDecoder(r.Body).Decode(&idsString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	var ids []uuid.UUID
	for _, idString := range idsString {
		id, err := uuid.Parse(idString)
		if err != nil {
			JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
			return
		}
		ids = append(ids, id)
	}

	err = h.productService.DeleteProducts(ctx, ids)
	if err != nil {
		JSONError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h ProductHandler) UpdateProductsFields(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.UpdateProductsFields()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	var jsonBody map[string]any
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	product, err := h.productService.UpdateProductFields(ctx, id, jsonBody)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: product,
			Meta: map[string]any{
				"fieldsUpdated": map[string]any{
					"total": len(jsonBody),
				},
			},
		},
	}, http.StatusOK, w)
}

func (h ProductHandler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.DeleteProductById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	err = h.productService.DeleteProductById(ctx, id)
	if err != nil {
		JSONError(w, err)
	}
	w.WriteHeader(http.StatusNoContent)
}
