package repository

import "calendar.com/pkg/storage"

type Repository struct {
	Storage *storage.Storage
	UserRepository
}

func NewRepository(storage *storage.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}
