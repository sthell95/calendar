package repository

import "calendar.com/pkg/domain/entity"

type EventRepository interface {
	Create(*entity.Event) error
	Update(event *entity.Event, id string) (*entity.Event, error)
	FindOneById(string) (*entity.Event, error)
}

type Event struct {
	repo Repository
}

func (e *Event) Create(event *entity.Event) error {
	return nil
}

func (e *Event) Update(event *entity.Event, id string) (*entity.Event, error) {
	return nil, nil
}

func (e *Event) FindOneById(id string) (*entity.Event, error) {
	return nil, nil
}
