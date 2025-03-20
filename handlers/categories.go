package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"rest-api-example/entities"
	"rest-api-example/service"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) CategoryHandler {
	return CategoryHandler{
		categoryService: s,
	}
}

func (h CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	categories, err := h.categoryService.GetAllCategories(ctx, queryParams)
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

func (h CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.GetCategoryById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	category, err := h.categoryService.GetCategoryById(ctx, id)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: category,
			Meta: map[string]any{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusOK, w)
}

func (h CategoryHandler) GetCategoriesByIds(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.GetCategoriesByIds()"
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

	categories, err := h.categoryService.GetCategoriesByIds(ctx, ids)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: categories,
			Meta: map[string]any{
				"result": map[string]int{
					"total": len(categories),
				},
			},
		},
	}, http.StatusOK, w)
}

func (h CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.CreateCategory()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var category entities.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	category, err = h.categoryService.CreateCategory(ctx, category)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: category,
			Meta: map[string]any{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusCreated, w)
}

func (h CategoryHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.DeleteCategoryById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	err = h.categoryService.DeleteCategoryById(ctx, id)
	if err != nil {
		JSONError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h CategoryHandler) DeleteCategories(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.DeleteCategories()"
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

	err = h.categoryService.DeleteCategories(ctx, ids)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h CategoryHandler) UpdateCategoryFields(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.UpdateCategoryFields()"
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

	category, err := h.categoryService.UpdateCategoryFields(ctx, id, jsonBody)
	if err != nil {
		JSONError(w, err)
		return
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: category,
			Meta: map[string]any{
				"fieldsUpdated": map[string]any{
					"total": len(jsonBody),
				},
			},
		},
	}, http.StatusOK, w)
}

func (h CategoryHandler) GetAllProductsByCategory(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.GetAllProductsByCategory()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	category, err := h.categoryService.GetCategoryById(ctx, id)
	if err != nil {
		JSONError(w, err)
		return
	}

	products, err := h.categoryService.GetAllProductsByCategory(ctx, id)
	if err != nil {
		JSONError(w, err)
		return
	}

	type CategoryProducts struct {
		Category entities.Category  `json:"category"`
		Products []entities.Product `json:"products"`
	}
	categoryProducts := CategoryProducts{
		Category: category,
		Products: products,
	}

	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: categoryProducts,
			Meta: map[string]any{
				"products": map[string]any{
					"total": len(products),
				},
			},
		},
	}, http.StatusOK, w)
}

func getQueryInt(query url.Values, key string, defaultValue int) int {
	if value := query.Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
