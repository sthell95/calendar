package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
)

type Authorization interface {
	SignIn(w http.ResponseWriter, r *http.Request)
}

type Services struct {
	service.AuthService
}

func (c Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		return
	}
	var credential entity.Credentials

	err = json.Unmarshal(data, &credential)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		return
	}

	_, err = c.Services.Authorization.UserRepository.FindByCredentials(credential)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		return
	}

	tokenString := jwt.SigningMethodHS256.Hash.String()
	err = c.Services.Authorization.GenerateJWT(&tokenString)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		return
	}
}
