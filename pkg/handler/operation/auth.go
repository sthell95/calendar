package operation

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/service"
)

type Auth struct {
	service.Authorization
}

func NewAuthOperations(s service.Authorization) *Auth {
	return &Auth{s}
}

func (c *Auth) SignIn(ctx context.Context, data *entity.Credentials) (*entity.AuthToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sign-in")
	defer span.Finish()

	token, err := c.Authorization.SignInProcess(ctx, data)
	if err != nil {
		return nil, err
	}

	return token, nil
}
