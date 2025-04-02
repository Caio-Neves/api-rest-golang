package category

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"strconv"
	"strings"
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
	op := "CategoryHandler.GetAllCategories()"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	page := utils.GetQueryInt(queryParams, "page", 1)
	limit := utils.GetQueryInt(queryParams, "limit", 10)

	categories, totalCount, err := h.categoryService.GetAllCategories(ctx, page, limit, queryParams)
	if err != nil {
		utils.JSONError(w, err)
		return
	}

	resources := make([]entities.CategoryResource, len(categories))
	for index, category := range categories {
		url := fmt.Sprintf("%s/%s", r.URL.Path, category.Id.String())
		links := entities.NewHateoasBuilder().
			AddGet("self", url).
			AddDelete("delete", fmt.Sprintf("/admin%s", url)).
			AddPatch("update", fmt.Sprintf("/admin%s", url)).
			Build()

		resource := entities.CategoryResource{
			Category: category,
			Links:    links,
		}
		resources[index] = resource
	}

	var filtersUrl string

	var isActive int
	if value, exists := queryParams["active"]; exists {
		isActive, err = strconv.Atoi(value[0])
		if err != nil {
			utils.JSONError(w, entities.NewInternalServerErrorError(err, op))
		}
		filtersUrl += fmt.Sprintf("&active=%d", isActive)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	paginationLinksBuilder := entities.NewHateoasBuilder().
		AddGet("self", fmt.Sprintf("%s?page=%d&limit=%d%s", r.URL.Path, page, limit, filtersUrl))
	if page < totalPages {
		paginationLinksBuilder.AddGet("last", fmt.Sprintf("%s?page=%d&limit=%d%s", r.URL.Path, totalPages, limit, filtersUrl))
	}
	if page+1 <= totalPages {
		paginationLinksBuilder.AddGet("next", fmt.Sprintf("%s?page=%d&limit=%d%s", r.URL.Path, page+1, limit, filtersUrl))
	}
	if page-1 > 0 {
		paginationLinksBuilder.AddGet("prev", fmt.Sprintf("%s?page=%d&limit=%d%s", r.URL.Path, page-1, limit, filtersUrl))
	}
	links := paginationLinksBuilder.Build()

	meta := utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Results:    len(categories),
		Hateoas:    links,
	}

	utils.JSONResponse(w, resources, meta, http.StatusOK)
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

	var url string
	url = strings.ReplaceAll(r.URL.Path, "\"", url)
	links := entities.NewHateoasBuilder().
		AddGet("self", url).
		AddDelete("delete", fmt.Sprintf("/admin%s", url)).
		AddPatch("update", fmt.Sprintf("/admin%s", url)).
		Build()

	utils.JSONResponse(w, category, links, http.StatusOK)
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

	resources := make([]entities.CategoryResource, len(categories))
	for index, category := range categories {
		url := fmt.Sprintf("/categories/%s", category.Id.String())
		links := entities.NewHateoasBuilder().
			AddGet("self", url).
			AddDelete("delete", fmt.Sprintf("/admin%s", url)).
			AddPatch("update", fmt.Sprintf("/admin%s", url)).
			Build()
		resource := entities.CategoryResource{
			Category: category,
			Links:    links,
		}
		resources[index] = resource
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

	url := fmt.Sprintf("/categories/%s", category.Id.String())
	links := entities.NewHateoasBuilder().
		AddGet("self", url).
		AddDelete("delete", fmt.Sprintf("/admin%s", url)).
		AddPatch("update", fmt.Sprintf("/admin%s", url)).
		Build()

	utils.JSONResponse(w, category, links, http.StatusCreated)
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

	url := fmt.Sprintf("/categories/%s", category.Id.String())
	links := entities.NewHateoasBuilder().
		AddGet("self", url).
		AddDelete("delete", fmt.Sprintf("/admin%s", url)).
		AddPatch("update", fmt.Sprintf("/admin%s", url)).
		Build()

	utils.JSONResponse(w, category, links, http.StatusOK)
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

	resources := make([]entities.ProductResource, len(products))
	for index, product := range products {
		url := fmt.Sprintf("%s/%s", "/products", product.Id.String())
		links := entities.NewHateoasBuilder().
			AddGet("self", url).
			AddDelete("delete", fmt.Sprintf("/admin%s", url)).
			AddPatch("update", fmt.Sprintf("/admin%s", url)).
			Build()
		resource := entities.ProductResource{
			Product: product,
			Links:   links,
		}
		resources[index] = resource
	}

	utils.JSONResponse(w, resources, nil, http.StatusOK)
}
