package product

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ErrIdDosProdutosObrigatorio = errors.New("id dos produtos a serem excluídos devem ser informados")
)

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(s ProductService) ProductHandler {
	return ProductHandler{
		productService: s,
	}
}

func (h ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	page := utils.GetQueryInt(queryParams, "page", 1)
	limit := utils.GetQueryInt(queryParams, "limit", 10)
	products, err := h.productService.GetAllProducts(ctx, queryParams)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	paginationMeta := utils.PaginationMeta{
		Page:    page,
		Limit:   limit,
		Results: len(products),
	}
	utils.JSONResponse(w, products, paginationMeta, http.StatusOK)
}

func (h ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.GetProductById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	product, err := h.productService.GetProductById(ctx, id)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, product, nil, http.StatusOK)
}

func (h ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.CreateProduct()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var product entities.Product
	log.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	_, err = h.productService.CreateProduct(ctx, product)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, product, nil, http.StatusCreated)
}

func (h ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.DeleteProducts()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var idsString []string
	err := json.NewDecoder(r.Body).Decode(&idsString)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	var ids []uuid.UUID
	for _, idString := range idsString {
		id, err := uuid.Parse(idString)
		if err != nil {
			utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
			return
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		utils.JSONError(w, entities.NewBadRequestError(ErrIdDosProdutosObrigatorio, ErrIdDosProdutosObrigatorio.Error(), op))
		return
	}

	err = h.productService.DeleteProducts(ctx, ids)
	if err != nil {
		utils.JSONError(w, err)
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
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	var jsonBody map[string]any
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	product, err := h.productService.UpdateProductFields(ctx, id, jsonBody)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, product, nil, http.StatusOK)
}

func (h ProductHandler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
	op := "ProductHandler.DeleteProductById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	err = h.productService.DeleteProductById(ctx, id)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
