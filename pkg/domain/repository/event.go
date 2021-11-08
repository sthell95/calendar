package repository

import (
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/storage"
)

type EventRepository interface {
	Create(event *entity.Event) error
	Update(event *entity.Event, id string) (*entity.Event, error)
	FindOneById(string) (*entity.Event, error)
}

type Event struct {
	repo storage.Repository
}

func (e *Event) Create(event *entity.Event) error {
	return e.repo.Create(&event)
}

func (e *Event) Update(event *entity.Event, id string) (*entity.Event, error) {
	return nil, nil
}

func (e *Event) FindOneById(id string) (*entity.Event, error) {
	return nil, nil
}

func NewEventRepository(repo storage.Repository) *Event {
	return &Event{
		repo: repo,
	}
}
