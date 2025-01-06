package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, username, password string) (models.LoginResponse, error) {
	dbPass, err := s.storage.GetUserPassword(ctx, username)
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("failed to retrieve user password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
	if err != nil {
		log.Println(err)
		return models.LoginResponse{}, ErrUnauthorized
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTCustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	jwtString, err := t.SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("failed to sign JWT: %w", err)
	}

	return models.LoginResponse{
		JWT: jwtString,
		Exp: time.Now().Add(time.Hour * 24).Unix(),
	}, nil
}
