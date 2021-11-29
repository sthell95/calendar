package controller

import (
	"encoding/json"
	"net/http"

	"github.com/opentracing/opentracing-go"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type RequestCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "sign-in")
	defer span.Finish()

	var credential RequestCredentials
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
		return
	}

	data := entity.Credentials{
		Login:    credential.Login,
		Password: credential.Password,
	}
	token, err := c.AuthService.SignInProcess(ctx, &data)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
		return
	}
	responseToken := ResponseToken{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	}
	response.NewPrint().PrettyPrint(w, responseToken)
}
