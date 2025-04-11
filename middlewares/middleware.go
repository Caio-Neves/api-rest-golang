package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"slices"
	"strings"
)

func ValidateSupportedMediaTypes(mediaTypes []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		op := "middlewares.ValidateSupportedMediaTypes()"
		contentType := strings.TrimSpace(r.Header.Get("Content-Type"))
		if contentType != "" {
			for _, mediaType := range mediaTypes {
				if strings.HasPrefix(contentType, mediaType) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}
		utils.JSONError(w, entities.NewUnsupportedMediaType(errors.New("formato não suportado"),
			fmt.Sprintf("tente os seguintes media types %s", strings.Join(mediaTypes, ",")), op))
	})
}

func ValidadeAcceptHeader(acceptContents []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		op := "middlewares.ValidadeAcceptHeader()"
		acceptRequest := strings.TrimSpace(r.Header.Get("Accept"))
		if slices.Contains(acceptContents, acceptRequest) || acceptRequest == "" || acceptRequest == "*/*" {
			next.ServeHTTP(w, r)
			return
		}
		utils.JSONError(w, entities.NewNotAcceptable(errors.New("formato não suportado"), fmt.Sprintf("formatos de retorno: %s", strings.Join(acceptContents, ",")), op))
	})
}
