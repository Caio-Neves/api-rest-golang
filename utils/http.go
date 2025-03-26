package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"rest-api-example/entities"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

type PaginationMeta struct {
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	Results int `json:"results"`
}

func GetQueryInt(query url.Values, key string, defaultValue int) int {
	if value := query.Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func JSONError(w http.ResponseWriter, err error) {
	e, ok := err.(*entities.Error)
	if !ok {
		e = entities.NewInternalServerErrorError(err, "Unhandled error")
	}

	switch e.Code {
	case entities.BAD_REQUEST:
		w.WriteHeader(http.StatusBadRequest)
	case entities.INTERNAL_SERVER_ERROR:
		w.WriteHeader(http.StatusInternalServerError)
	case entities.NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
	case entities.CONFLICT:
		w.WriteHeader(http.StatusConflict)
	case entities.UNAUTHORIZED:
		w.WriteHeader(http.StatusUnauthorized)
	case entities.FORBIDDEN:
		w.WriteHeader(http.StatusForbidden)
	case entities.NOT_IMPLEMENTED:
		w.WriteHeader(http.StatusNotImplemented)
	case entities.UNSUPPORTED_MEDIA_TYPE:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case entities.NOT_ACCEPTABLE:
		w.WriteHeader(http.StatusNotAcceptable)
	}

	entry := log.WithFields(log.Fields{
		"code":      e.Code,
		"error":     e.Err.Error(),
		"operation": e.Operation,
		"message":   e.Message,
	})

	if e.Code != entities.INTERNAL_SERVER_ERROR {
		entry.Info()
	} else {
		entry.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func JSONResponse(w http.ResponseWriter, data any, meta any, statusCode int) {
	op := "utils.JSONResponse()"
	r := Response{
		Data: data,
		Meta: meta,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		JSONError(w, entities.NewInternalServerErrorError(err, op))
		return
	}
}
