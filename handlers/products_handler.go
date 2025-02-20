package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"rest-api-example/entities"
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
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var product entities.Product
	log.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		SendJsonError(JsonResponseError{
			Payload: ResponseError{
				Error: "Verifique o formato do JSON e tente novamente.",
			},
		}, http.StatusBadRequest, w)
		return
	}
	product, err = h.productService.CreateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, errorsApi.ErrCategoriaDoProdutoEhObrigatoria) ||
			errors.Is(err, errorsApi.ErrCategoriaNaoCadastrada) ||
			errors.Is(err, errorsApi.ErrDescricaoProdutoEhObrigatorio) ||
			errors.Is(err, errorsApi.ErrNomeProdutoEhObrigatorio) {
			SendJsonError(JsonResponseError{
				Payload: ResponseError{
					Error: err.Error(),
				},
			}, http.StatusBadRequest, w)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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
	}, http.StatusCreated, w)
}

func (h ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var idsString []string
	err := json.NewDecoder(r.Body).Decode(&idsString)
	if err != nil {
		SendJsonError(JsonResponseError{
			Payload: ResponseError{
				Error: "Verifique o formato do JSON e tente novamente.",
			},
		}, http.StatusBadRequest, w)
		return
	}
	var ids []uuid.UUID
	for _, idString := range idsString {
		id, err := uuid.Parse(idString)
		if err != nil {
			SendJsonError(JsonResponseError{
				Payload: ResponseError{
					Error: errorsApi.ErrUuidInvalido.Error(),
				},
			}, http.StatusBadRequest, w)
			return
		}
		ids = append(ids, id)
	}
	err = h.productService.DeleteProducts(ctx, ids)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h ProductHandler) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: UpdateProducts"))
}

func (h ProductHandler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
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
	err = h.productService.DeleteProductById(ctx, id)
	if err != nil {
		switch err {
		case errorsApi.ErrProdutoNaoCdastrado:
			SendJsonError(JsonResponseError{
				Payload: ResponseError{
					Error: "Produto não cadastrado.",
				},
			}, http.StatusBadRequest, w)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
