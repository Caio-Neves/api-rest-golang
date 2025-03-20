package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api-example/entities"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

type JsonResponse struct {
	Header  http.Header
	Payload Response
}

type ResponseError struct {
	Error               interface{} `json:"error"`
	SupportedMediaTypes []string    `json:"supported_formats,omitempty"`
	UnknownFields       []string    `json:"unknown_fields,omitempty"`
}

type JsonResponseError struct {
	Header  http.Header
	Payload ResponseError
}

func SendJsonResponse(jsonResponse JsonResponse, httpStatusCode int, w http.ResponseWriter) {
	for key, values := range jsonResponse.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(w).Encode(jsonResponse.Payload)
	if err != nil {
		return
	}
}

func SendJsonError(jsonResponseError JsonResponseError, httpStatusCode int, w http.ResponseWriter) {
	for key, values := range jsonResponseError.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(w).Encode(jsonResponseError.Payload)
	if err != nil {
		return
	}
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

	json.NewEncoder(w).Encode(err)
}
