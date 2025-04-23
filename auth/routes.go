package auth

import (
	"net/http"
	"rest-api-example/middlewares"
	"rest-api-example/user"

	"github.com/gorilla/mux"
)

func SetupAuthRoutes(mux *mux.Router, authHandler AuthHandler, userHandler user.UserHandler) {
	authRoutes := mux.PathPrefix("/auth").Subrouter()
	authRoutes.Path("/login").HandlerFunc(
		middlewares.ValidadeAcceptHeader([]string{"application/json"}, authHandler.Login)).Methods(http.MethodPost)
	authRoutes.Path("/refresh").HandlerFunc(
		middlewares.ValidadeAcceptHeader([]string{"application/json"}, authHandler.RefreshToken)).Methods(http.MethodPost)
}
