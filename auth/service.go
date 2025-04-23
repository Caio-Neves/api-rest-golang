package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rest-api-example/entities"
	"rest-api-example/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpectedAccessToken  = errors.New("expected access token")
	ErrExpectedRefreshToken = errors.New("expected refresh token")
)

type AccessToken string
type RefreshToken string

type TokenPair struct {
	AccessToken
	RefreshToken
}

type AuthService struct {
	userRepository entities.UserInterface
	secretKey      string
}

func NewAuthService(userRepository entities.UserInterface, secretKey string) AuthService {
	return AuthService{
		userRepository: userRepository,
		secretKey:      secretKey,
	}
}

func (u AuthService) validateToken(tokenString string) (*jwt.Token, error) {
	op := "AuthService.ValidateToken()"

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(u.secretKey), nil
	})
	if err != nil {
		return nil, entities.NewUnauthorizedError(err, err.Error(), op)
	}

	if !token.Valid {
		return nil, entities.NewUnauthorizedError(ErrInvalidToken, ErrInvalidToken.Error(), op)
	}
	return token, nil
}

func (u AuthService) Login(ctx context.Context, credentials entities.Credentials) (TokenPair, error) {
	op := "AuthService.Login()"
	credentialsDatabase, err := u.userRepository.GetCredentialsByLogin(ctx, credentials.Login)
	if err != nil {
		return TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentialsDatabase.Password), []byte(credentials.Password))
	if err != nil {
		return TokenPair{}, entities.NewUnauthorizedError(ErrInvalidCredentials, ErrInvalidCredentials.Error(), op)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  credentials.Login,
		"iss":  "ecomapi",
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"iat":  time.Now().Unix(),
		"type": "access_token",
	})

	signedAccessToken, err := accessToken.SignedString([]byte(u.secretKey))
	if err != nil {
		return TokenPair{}, entities.NewInternalServerErrorError(err, op)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  credentials.Login,
		"iss":  "ecomapi",
		"exp":  time.Now().Add(time.Hour * (24 * 7)).Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh_token",
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(u.secretKey))
	if err != nil {
		return TokenPair{}, entities.NewInternalServerErrorError(err, op)
	}

	return TokenPair{AccessToken(signedAccessToken), RefreshToken(signedRefreshToken)}, nil
}

func (u AuthService) RefreshToken(refreshToken RefreshToken) (TokenPair, error) {
	op := "AuthService.RefreshToken()"

	token, err := u.validateToken(string(refreshToken))
	if err != nil {
		return TokenPair{}, entities.NewUnauthorizedError(err, err.Error(), op)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenPair{}, entities.NewUnauthorizedError(errors.New("error parse token claims"), "error parse token claims", op)
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return TokenPair{}, entities.NewUnauthorizedError(errors.New("subject not found on token"), "subject not found on token", op)
	}

	typeToken, ok := claims["type"].(string)
	if !ok {
		return TokenPair{}, entities.NewInternalServerErrorError(errors.New("error parse token claims"), op)
	}

	if typeToken != "refresh_token" {
		return TokenPair{}, entities.NewUnauthorizedError(ErrExpectedRefreshToken, ErrExpectedRefreshToken.Error(), op)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  sub,
		"iss":  "ecomapi",
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"iat":  time.Now().Unix(),
		"type": "access_token",
	})

	signedAccessToken, err := accessToken.SignedString([]byte(u.secretKey))
	if err != nil {
		return TokenPair{}, entities.NewInternalServerErrorError(err, op)
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  sub,
		"iss":  "ecomapi",
		"exp":  time.Now().Add(time.Hour * (24 * 7)).Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh_token",
	})

	signedRefreshToken, err := newRefreshToken.SignedString([]byte(u.secretKey))
	if err != nil {
		return TokenPair{}, entities.NewInternalServerErrorError(err, op)
	}

	return TokenPair{AccessToken(signedAccessToken), RefreshToken(signedRefreshToken)}, nil
}

func (u AuthService) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		op := "AuthService.AuthenticationMiddleware()"
		tokenString, err := utils.GetBearerToken(r)
		if err != nil {
			utils.JSONError(w, err)
			return
		}

		token, err := u.validateToken(tokenString)
		if err != nil {
			utils.JSONError(w, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.JSONError(w, entities.NewInternalServerErrorError(errors.New("error parse token claims"), op))
			return
		}

		typeToken, ok := claims["type"].(string)
		if !ok {
			utils.JSONError(w, entities.NewInternalServerErrorError(errors.New("error parse token claims"), op))
			return
		}

		if typeToken != "access_token" {
			utils.JSONError(w, entities.NewUnauthorizedError(ErrExpectedAccessToken, ErrExpectedAccessToken.Error(), op))
			return
		}

		log.Printf("Token claims %+v\\n", token.Claims)
		next.ServeHTTP(w, r)
	})
}
