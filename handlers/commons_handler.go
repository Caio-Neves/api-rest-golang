package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	errorsApi "rest-api-example/errors"
	"strings"
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

func ValidateJSONFields(jsonMap map[string]interface{}, genStruct interface{}) ([]string, error) {
	structVal := reflect.ValueOf(genStruct)
	if structVal.Kind() != reflect.Ptr || structVal.IsNil() {
		return nil, errors.New("item deve ser um ponteiro pra uma struct")
	}
	structVal = structVal.Elem()
	if structVal.Kind() != reflect.Struct {
		return nil, errors.New("item não é uma struct")
	}
	structType := structVal.Type()
	var unknownFields []string
	for key := range jsonMap {
		fieldFound := false
		for i := 0; i < structVal.NumField(); i++ {
			tag := structType.Field(i).Tag.Get("json")
			if strings.EqualFold(key, strings.SplitN(tag, ",", 2)[0]) {
				fieldFound = true
				break
			}
		}
		if !fieldFound {
			unknownFields = append(unknownFields, key)
		}
	}
	if len(unknownFields) > 0 {
		return unknownFields, errorsApi.ErrAtributoNaoExistente
	}
	return nil, nil
}
