package user

import (
	"net/http"
	"rest-api-example/middlewares"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(mux *mux.Router, h UserHandler) {
	userRoutes := mux.PathPrefix("/users").Subrouter()
	userRoutes.Path("").HandlerFunc(
		middlewares.ValidadeAcceptHeader([]string{"application/json"}, h.RegisterUser)).Methods(http.MethodPost)
}
