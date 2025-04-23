package entities

import "context"

type UserInterface interface {
	GetCredentialsByLogin(ctx context.Context, login string) (Credentials, error)
	InsertUser(ctx context.Context, credentials Credentials) error
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
