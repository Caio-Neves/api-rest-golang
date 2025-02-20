package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"rest-api-example/entities"
	errorsApi "rest-api-example/errors"
	"rest-api-example/service"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: s,
	}
}

func (h CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	queryParams := r.URL.Query()
	categories, err := h.categoryService.GetAllCategories(ctx, queryParams)
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

func (h CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
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
	category, err := h.categoryService.GetCategoryById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if category.IsEmpty() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: category,
			Meta: map[string]interface{}{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusOK, w)
}

func (h CategoryHandler) GetCategoriesByIds(w http.ResponseWriter, r *http.Request) {
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
	categories, err := h.categoryService.GetCategoriesByIds(ctx, ids)
	if err != nil {
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
				"result": map[string]int{
					"total": len(categories),
				},
			},
		},
	}, http.StatusOK, w)
}

func (h CategoryHandler) GetAllProductsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint: GetAllProductsByCategory"))
}

func (h CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var category entities.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		SendJsonError(JsonResponseError{
			Payload: ResponseError{
				Error: "Verifique o formato do JSON e tente novamente.",
			},
		}, http.StatusBadRequest, w)
		return
	}
	category, err = h.categoryService.CreateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, errorsApi.ErrNomeCategoriaObrigatorio) || errors.Is(err, errorsApi.ErrDescricaoCategoriaObrigatorio) {
			SendJsonError(JsonResponseError{
				Payload: ResponseError{
					Error: err.Error(),
				},
			}, http.StatusBadRequest, w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: category,
			Meta: map[string]interface{}{
				"result": map[string]int{
					"total": 1,
				},
			},
		},
	}, http.StatusCreated, w)
}

func (h CategoryHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
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
	err = h.categoryService.DeleteCategoryById(ctx, id)
	if err != nil {
		switch err {
		case errorsApi.ErrCategoriaNaoCadastrada:
			SendJsonError(JsonResponseError{
				Payload: ResponseError{
					Error: "Categoria não cadastrada.",
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

func (h CategoryHandler) DeleteCategories(w http.ResponseWriter, r *http.Request) {
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
	err = h.categoryService.DeleteCategories(ctx, ids)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h CategoryHandler) UpdateCategoryFields(w http.ResponseWriter, r *http.Request) {
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
	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		SendJsonError(JsonResponseError{
			Payload: ResponseError{
				Error: "Verifique o formato do JSON e tente novamente.",
			},
		}, http.StatusBadRequest, w)
		return
	}
	category := entities.Category{}
	unknownFields, err := ValidateJSONFields(jsonBody, &category)
	if err != nil {
		if errors.Is(err, errorsApi.ErrAtributoNaoExistente) {
			SendJsonError(JsonResponseError{
				Payload: ResponseError{
					Error:         err.Error(),
					UnknownFields: unknownFields,
				},
			}, http.StatusBadRequest, w)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	category, err = h.categoryService.UpdateCategoryFields(ctx, id, jsonBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SendJsonResponse(JsonResponse{
		Payload: Response{
			Data: category,
			Meta: map[string]interface{}{
				"fieldsUpdated": map[string]interface{}{
					"total": len(jsonBody),
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
