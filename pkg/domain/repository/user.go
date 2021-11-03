package repository

import (
	"github.com/gofrs/uuid"

	"calendar.com/pkg/storage"

	"calendar.com/pkg/domain/entity"
)

type UserRepository interface {
	Create(*entity.User) error
	FindById(uuid.UUID) (*entity.User, error)
	FindOneBy(conditions map[string]interface{}) (*entity.User, error)
}

type UserRepo struct {
	repo storage.Repository
}

func (r *UserRepo) Create(user *entity.User) error {
	return r.repo.Create(&user)
}

func (r *UserRepo) FindById(id uuid.UUID) (*entity.User, error) {
	var u entity.User
	err := r.repo.FindById(&u, id)

	return &u, err
}

func (r *UserRepo) FindOneBy(conditions map[string]interface{}) (*entity.User, error) {
	var u entity.User
	if err := r.repo.FindOneBy(&u, conditions); err != nil {
		return nil, err
	}
	return &u, nil
}

func NewUserRepository(r storage.Repository) *UserRepo {
	return &UserRepo{repo: r}
}
