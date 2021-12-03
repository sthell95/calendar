package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"calendar.com/pkg/domain/entity"
)

type UserRepository interface {
	FindById(context.Context, uuid.UUID) (*entity.User, error)
	FindOneByLogin(context.Context, string) (*entity.User, error)
}

type UserRepo struct {
	Repos []UserRepository
}

func NewUserRepository(repos ...UserRepository) *UserRepo {
	return &UserRepo{Repos: repos}
}

func (r *UserRepo) FindById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	for i, repo := range r.Repos {
		user, err := repo.FindById(ctx, id)
		if user == nil && i != len(r.Repos)-1 {
			continue
		}

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, errors.New("no providers")
}

func (r *UserRepo) FindOneByLogin(ctx context.Context, login string) (*entity.User, error) {
	for i, repo := range r.Repos {
		user, err := repo.FindOneByLogin(ctx, login)
		if user == nil && i != len(r.Repos)-1 {
			continue
		}

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, errors.New("no providers")
}
