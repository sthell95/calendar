package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type Error struct {
	Message string `json:"message"`
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		return
	}
	var credential entity.Credentials

	err = json.Unmarshal(data, &credential)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
		return
	}

	if err = c.Services.Authorization.CheckCredentials(credential); err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
		return
	}

	token, err := c.Services.Authorization.GenerateToken(credential.Login)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
		return
	}
	response.NewPrint().PrettyPrint(w, token)
}
