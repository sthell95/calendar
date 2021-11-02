package repository

import (
	"github.com/gofrs/uuid"

	"calendar.com/pkg/storage"

	"calendar.com/pkg/domain/entity"
)

type UserRepository interface {
	Create(*entity.User) error
	FindById(uuid.UUID) (*entity.User, error)
	FindByCredentials(entity.Credentials) (*entity.User, error)
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

func (r UserRepo) FindByCredentials(c entity.Credentials) (*entity.User, error) {
	u := entity.User{Login: c.Login}
	err := r.repo.Find(&u)

	return &u, err
}

func NewUserRepository(r storage.Repository) *UserRepo {
	return &UserRepo{repo: r}
}
