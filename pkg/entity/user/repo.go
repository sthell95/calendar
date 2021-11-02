package user

import (
	"calendar.com/pkg/storage/postgres"
)

type UserRepository interface {
	Create(user *User) error
}

type UserRepo struct {
	repo postgres.Repository
}

func (r *UserRepo) Create(user *User) error {
	return r.repo.Create(&user)
}

func NewUserRepo(repository postgres.Repository) *UserRepo {
	return &UserRepo{repo: repository}
}
