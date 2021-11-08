package repository

import (
	"calendar.com/pkg/storage"
)

type Repository struct {
	Storage *storage.Storage
	Event   EventRepository
}

func NewRepository(storage *storage.Storage) *Repository {
	return &Repository{
		Storage: storage,
		Event:   NewEventRepository(storage),
	}
}
