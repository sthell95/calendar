package repository

import (
	"github.com/gofrs/uuid"

	"calendar.com/pkg/storage"

	"calendar.com/pkg/domain/entity"
)

var _ UserRepository = (*User)(nil)

type UserRepository interface {
	FindById(uuid.UUID) (*entity.User, error)
	FindOneBy(conditions map[string]interface{}) (*entity.User, error)
}

type User struct {
	repo storage.Repository
}

func (r *User) FindById(id uuid.UUID) (*entity.User, error) {
	var u entity.User
	err := r.repo.FindById(&u, id)

	return &u, err
}

func (r *User) FindOneBy(conditions map[string]interface{}) (*entity.User, error) {
	var u entity.User
	if err := r.repo.FindOneBy(&u, conditions); err != nil {
		return nil, err
	}
	return &u, nil
}

func NewUserRepository(r storage.Repository) *User {
	return &User{repo: r}
}
