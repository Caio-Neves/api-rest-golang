package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"time"
)

var (
	ErrInvalidJsonFormat = errors.New("invalid json format")
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return UserHandler{userService: userService}
}

func (h UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.RegisterUser()"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var credentials entities.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, ErrInvalidJsonFormat.Error(), op))
		return
	}

	err = h.userService.Registry(ctx, credentials)
	if err != nil {
		utils.JSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
