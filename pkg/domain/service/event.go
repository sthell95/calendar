package service

import (
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
)

type Event interface {
	Create(event *entity.Event) error
}

type EventService struct {
	Repository repository.EventRepository
}

func (es *EventService) Create(e *entity.Event) error {
	return es.Repository.Create(e)
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{
		Repository: repo,
	}
}
