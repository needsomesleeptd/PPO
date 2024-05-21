package auth_utils

import (
	"annotater/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	errors2 "github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
)

type Payload struct {
	Login string
	ID    uint64
	Role  models.Role
}

type ITokenHandler interface {
	GenerateToken(credentials models.User, key string) (string, error)
	ValidateToken(tokenString string, key string) error
	ParseToken(tokenString string, key string) (*Payload, error)
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrParsingToken = errors.New("error parsing token")
)

type JWTTokenHandler struct {
}

func NewJWTTokenHandler() ITokenHandler {
	return JWTTokenHandler{}
}

func (hasher JWTTokenHandler) GenerateToken(credentials models.User, key string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exprires": time.Now().Add(time.Hour * 24),
			"login":    credentials.Login,
			"ID":       credentials.ID,
			"Role":     credentials.Role,
		})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("creating token err: %w", err)
	}

	return tokenString, nil
}

func (hasher JWTTokenHandler) ValidateToken(tokenString string, key string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return ErrParsingToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (hasher JWTTokenHandler) ParseToken(tokenString string, key string) (*Payload, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, errors2.Wrap(err, "auth.tokenhelper.GetRole error in parse")
	}

	payload := &Payload{
		Login: claims["login"].(string),
		ID:    uint64(claims["ID"].(float64)),
		Role:  models.Role(claims["Role"].(float64)),
	}

	return payload, nil
}
func ExtractTokenFromReq(r *http.Request) string {

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}
