package grpc

import (
	"context"

	"calendar.com/pkg/domain/entity"
	pg "calendar.com/proto"
)

type AuthOperations interface {
	SignIn(ctx context.Context, credentials *entity.Credentials) (*entity.AuthToken, error)
}

type AuthHandler struct {
	pg.UnimplementedAuthServiceServer
	AuthOperations AuthOperations
}

func (h *AuthHandler) Login(ctx context.Context, c *pg.Credentials) (*pg.Token, error) {
	entityCredentials := credentialsRequestToDomain(c)

	e, err := h.AuthOperations.SignIn(ctx, entityCredentials)
	if err != nil {
		return nil, err
	}

	return domainToResponseToken(e), nil
}

func credentialsRequestToDomain(c *pg.Credentials) *entity.Credentials {
	return &entity.Credentials{
		Login:    c.Login,
		Password: c.Password,
	}
}

func domainToResponseToken(t *entity.AuthToken) *pg.Token {
	return &pg.Token{
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
	}
}
