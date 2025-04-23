package user

import (
	"context"
	"rest-api-example/entities"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository entities.UserInterface
}

func NewUserService(u entities.UserInterface) UserService {
	return UserService{u}
}

func (u UserService) Registry(ctx context.Context, credentials entities.Credentials) error {
	op := "UserService.Registry()"
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.NewInternalServerErrorError(err, op)
	}
	credentials.Password = string(hashedPass)
	return u.userRepository.InsertUser(ctx, credentials)
}
