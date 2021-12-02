package postgresdb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"calendar.com/pkg/domain/entity"
)

var _ UserRepository = (*Client)(nil)

const userTable = "users"

type UserRepository interface {
	FindById(uuid.UUID) (*entity.User, error)
	FindOneBy(conditions map[string]interface{}) (*entity.User, error)
}

type user struct {
	ID       uuid.UUID
	Login    string
	Password string
	Timezone string
}

type UserModel struct{}

type Client struct {
	db *gorm.DB
}

func (r *Client) FindById(id uuid.UUID) (*entity.User, error) {
	var u entity.User
	err := r.db.Table(userTable).First(u, &id).Error

	return &u, err
}

func (r *Client) FindOneBy(conditions map[string]interface{}) (*entity.User, error) {
	var u entity.User
	err := r.db.Table(userTable).Take(&u, conditions).Error

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
