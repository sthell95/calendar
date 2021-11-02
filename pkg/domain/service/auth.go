package service

import (
	"crypto/sha1"
	"time"

	"calendar.com/pkg/logger"

	"github.com/golang-jwt/jwt"
)

type Auth struct{}

type Authorization interface {
	GenerateJWT(tokenString *string) error
	GeneratePassword(password string) string
}

const salt = "weg2c3928ncy29v2o3c23r29n3"

func (Auth) GenerateJWT(tokenString *string) error {
	token, err := jwt.ParseWithClaims(*tokenString, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}, nil)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "token")
		return err
	}
	ts, err := token.SigningString()
	tokenString = &ts

	return err
}

func (Auth) GeneratePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Sum([]byte(salt))

	return "qwe"
}
