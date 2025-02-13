package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
	UnknownFields       []string    `json:"unknown_fields"`
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
		log.Println(err)
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
		log.Println(err)
		return
	}
}
