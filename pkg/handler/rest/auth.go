package rest

import (
	"context"
	"encoding/json"
	"net/http"

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

type AuthOperations interface {
	SignIn(ctx context.Context, e *entity.Credentials) (*entity.AuthToken, error)
}

func (c *Client) SignIn(w http.ResponseWriter, r *http.Request) {
	var rc RequestCredentials
	err := json.NewDecoder(r.Body).Decode(&rc)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}

	entityCredentials := rc.credentialsRequestToDomain()
	token, err := c.AuthOperations.SignIn(r.Context(), entityCredentials)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}

	resp := domainToTokenResponse(token)
	w.WriteHeader(http.StatusOK)
	response.NewPrint().PrettyPrint(w, resp)
}

func (r *RequestCredentials) credentialsRequestToDomain() *entity.Credentials {
	return &entity.Credentials{
		Login:    r.Login,
		Password: r.Password,
	}
}

func domainToTokenResponse(e *entity.AuthToken) *ResponseToken {
	return &ResponseToken{
		Token:     e.Token,
		ExpiresAt: e.ExpiresAt,
	}
}
