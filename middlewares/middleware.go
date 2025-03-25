package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"strings"
)

func ValidateSupportedMediaTypes(mediaTypes []string, next http.HandlerFunc) http.HandlerFunc {
	op := "middlewares.ValidateSupportedMediaTypes()"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := strings.TrimSpace(r.Header.Get("Content-Type"))
		for _, mediaType := range mediaTypes {
			if contentType != "" && strings.HasPrefix(contentType, mediaType) {
				next.ServeHTTP(w, r)
				return
			}
		}
		utils.JSONError(w, entities.NewUnsupportedMediaType(errors.New("formato não suportado"), fmt.Sprintf("tente os seguintes media types %s", strings.Join(mediaTypes, ",")), op))
	})
}
