package handler

import (
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/response"
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"io"
)

type RequestCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type Auth struct {
	service.Authorization
}

func (c *Auth) SignIn(ctx context.Context, w io.Writer, r io.Reader) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sign-in")
	defer span.Finish()

	var credential RequestCredentials
	err := json.NewDecoder(r).Decode(&credential)
	if err != nil {
		return err
	}

	data := entity.Credentials{
		Login:    credential.Login,
		Password: credential.Password,
	}
	token, err := c.Authorization.SignInProcess(ctx, &data)
	if err != nil {
		return err
	}
	responseToken := ResponseToken{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	}

	response.NewPrint().PrettyPrint(w, responseToken)
	return nil
}
