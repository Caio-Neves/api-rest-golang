package service

import (
	"context"
	"rest-api-example/entities"
)

type UserService struct {
	userRepository entities.UserInterface
}

func NewUserService(u entities.UserInterface) UserService {
	return UserService{u}
}

func (u UserService) Login(ctx context.Context, credentials entities.Credentials) (bool, error) {
	auth, err := u.userRepository.CheckUserCredentials(ctx, credentials)
	if err != nil {
		return false, err
	}
	return auth, nil
}

func (u UserService) Registry(ctx context.Context, credentials entities.Credentials) error {
	return u.userRepository.RegistryUser(ctx, credentials)
}
