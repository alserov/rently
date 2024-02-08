package service

import (
	"encoding/base64"
	"fmt"
	"github.com/alserov/rently/user/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

const (
	ENV_SECRET_KEY = "SECRET_KEY"
)

func newToken(uuid string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, models.Claims{
		UUID: uuid,
		Role: role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 21)),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv(ENV_SECRET_KEY)))
	if err != nil {
		return "", &models.Error{
			Msg:    fmt.Sprintf("failed to sign token: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return tokenString, err
}

func parseTokenClaims(token string) (string, string, error) {
	c := models.Claims{}

	_, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(ENV_SECRET_KEY)), nil
	})
	if err != nil {
		return "", "", &models.Error{
			Msg:    "invalid token provided",
			Status: http.StatusBadRequest,
		}
	}

	return c.UUID, c.Role, nil
}

func hash(value string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.Error{
			Msg:    fmt.Sprintf("failed to generate hash from value: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return string(b), nil
}

func compareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return &models.Error{
			Msg:    "invalid password",
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

func encrypt(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decrypt(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", &models.Error{
			Msg:    fmt.Sprintf("failed to decrypt value: %v", err),
			Status: http.StatusInternalServerError,
		}
	}
	return string(data), nil
}
