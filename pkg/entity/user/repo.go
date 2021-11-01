package user

import (
	"github.com/gofrs/uuid"

	"calendar.com/pkg/storage/postgres"
)

type UserRepository interface {
	Create(user *User)
}

type UserRepo struct {
	repo postgres.Repository
}

func (r *UserRepo) Create(user *User) {
	r.repo.Create(&user)
}

func (r *UserRepo) Delete(ID uuid.UUID) {
	r.repo.Delete(ID)
}

func NewUserRepo(repository postgres.Repository) *UserRepo {
	return &UserRepo{repo: repository}
}
