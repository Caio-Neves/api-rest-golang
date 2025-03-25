package category

import (
	"context"
	"encoding/json"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(s CategoryService) CategoryHandler {
	return CategoryHandler{
		categoryService: s,
	}
}

func (h CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	page := utils.GetQueryInt(queryParams, "page", 1)
	limit := utils.GetQueryInt(queryParams, "limit", 10)

	categories, err := h.categoryService.GetAllCategories(ctx, page, limit, queryParams)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	paginationMeta := utils.PaginationMeta{
		Page:    page,
		Limit:   limit,
		Results: len(categories),
	}
	utils.JSONResponse(w, categories, paginationMeta, http.StatusOK)
}

func (h CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.GetCategoryById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	category, err := h.categoryService.GetCategoryById(ctx, id)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, category, nil, http.StatusOK)
}

func (h CategoryHandler) GetCategoriesByIds(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.GetCategoriesByIds()"
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

	categories, err := h.categoryService.GetCategoriesByIds(ctx, ids)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, categories, nil, http.StatusOK)
}

func (h CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.CreateCategory()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var category entities.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	category, err = h.categoryService.CreateCategory(ctx, category)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, category, nil, http.StatusCreated)
}

func (h CategoryHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.DeleteCategoryById()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	err = h.categoryService.DeleteCategoryById(ctx, id)
	if err != nil {
		utils.JSONError(w, err)
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
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	var jsonBody map[string]any
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "Verifique o formato do JSON e tente novamente", op))
		return
	}

	category, err := h.categoryService.UpdateCategoryFields(ctx, id, jsonBody)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, category, nil, http.StatusOK)
}

func (h CategoryHandler) GetAllProductsByCategory(w http.ResponseWriter, r *http.Request) {
	op := "CategoryHandler.GetAllProductsByCategory()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := uuid.Parse(idString)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, "UUID inválido", op))
		return
	}

	_, err = h.categoryService.GetCategoryById(ctx, id)
	if err != nil {
		utils.JSONError(w, err)
		return
	}

	products, err := h.categoryService.GetAllProductsByCategory(ctx, id)
	if err != nil {
		utils.JSONError(w, err)
		return
	}
	utils.JSONResponse(w, products, nil, http.StatusOK)
}
