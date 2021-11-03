package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"calendar.com/pkg/response"

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

	if err := c.Services.Authorization.CheckCredentials(credential); err != nil {
		logger.NewLogger().Write(logger.Error, "Invalid credentials", "sign-in")
		return
	}

	token, err := c.Services.Authorization.GenerateToken(credential.Login)
	response.NewPrint().PrettyPrint(w, token)
}
