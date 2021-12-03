package postgresdb

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"calendar.com/pkg/domain/repository"

	"calendar.com/pkg/domain/entity"
)

var _ repository.UserRepository = (*Client)(nil)

const userTable = "users"

type user struct {
	ID       uuid.UUID
	Login    string
	Password string
	Timezone string
}

type UserRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindOneByLogin(ctx context.Context, login string) (*entity.User, error)
}

type UserModel struct{}

type Client struct {
	db *gorm.DB
}

func (r *Client) FindById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Find user in postgres")
	defer span.Finish()

	var u entity.User
	err := r.db.Table(userTable).First(u, &id).Error

	return &u, err
}

func (r *Client) FindOneByLogin(ctx context.Context, login string) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Find user in postgres")
	defer span.Finish()

	var u entity.User
	err := r.db.Table(userTable).Take(&u, struct{ login string }{login: login}).Error

	return &u, err
}

func (*UserModel) toModel(u *entity.User) *user {
	return &user{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
		Timezone: u.Timezone,
	}
}

func (*UserModel) toDomainModel(u *user) *entity.User {
	return &entity.User{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
		Timezone: u.Timezone,
	}
}
