package middlewares

import (
	"net/http"
	"rest-api-example/handlers"
)

func ValidateSupportedMediaTypes(mediaTypes []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, mediaType := range mediaTypes {
			if r.Header.Get("Content-Type") != "" && r.Header.Get("Content-Type") == mediaType {
				next.ServeHTTP(w, r)
				return
			}
		}
		handlers.SendJsonError(handlers.JsonResponseError{
			Payload: handlers.ResponseError{
				Error:               "Formato não suportado",
				SupportedMediaTypes: mediaTypes,
			},
		}, http.StatusUnsupportedMediaType, w)
	})
}
