package entities

import "context"

type UserInterface interface {
	CheckUserCredentials(ctx context.Context, credentials Credentials) (bool, error)
	RegistryUser(ctx context.Context, credentials Credentials) error
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
