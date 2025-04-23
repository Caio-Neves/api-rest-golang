package auth

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

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

type AuthenticationResponse struct {
	AccessToken  `json:"access_token"`
	RefreshToken `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken `json:"refresh_token"`
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.Login()"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var credentials entities.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, ErrInvalidJsonFormat.Error(), op))
		return
	}

	tokenPair, err := h.authService.Login(ctx, credentials)
	if err != nil {
		utils.JSONError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(AuthenticationResponse{AccessToken(tokenPair.AccessToken), RefreshToken(tokenPair.RefreshToken)})
	if err != nil {
		utils.JSONError(w, err)
		return
	}
}

func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.RefreshToken()"

	var refreshTokenRequest RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest)
	if err != nil {
		utils.JSONError(w, entities.NewBadRequestError(err, ErrInvalidJsonFormat.Error(), op))
		return
	}

	tokenPair, err := h.authService.RefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		utils.JSONError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(AuthenticationResponse{AccessToken(tokenPair.AccessToken), RefreshToken(tokenPair.RefreshToken)})
	if err != nil {
		utils.JSONError(w, err)
		return
	}
}
