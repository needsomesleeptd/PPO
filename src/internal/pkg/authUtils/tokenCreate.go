package auth_utils

import (
	"annotater/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ITokenHandler interface {
	GenerateToken(credentials models.User, key string) (string, error)
	ValidateToken(tokenString string, key string) error
}

type JWTTokenHandler struct {
}

func NewJWTTokenHandler() ITokenHandler {
	return JWTTokenHandler{}
}

func (hasher JWTTokenHandler) GenerateToken(credentials models.User, key string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exprires": time.Now().Add(time.Hour * 24).Unix(),
			"login":    credentials.Login,
			"role":     credentials.Role,
		})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("creating token err: %w", err)
	}

	return tokenString, nil
}

func (hasher JWTTokenHandler) ValidateToken(tokenString string, key string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return errors.New("token is invalid")
	}

	return nil
}
