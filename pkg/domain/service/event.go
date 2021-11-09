package service

import (
	"calendar.com/pkg/controller"
	"calendar.com/pkg/domain/repository"
)

type Event interface {
	Create(*controller.RequestEvent) error
}

type EventService struct {
	Repository repository.EventRepository
}

func (es *EventService) Create(e *controller.RequestEvent) error {

	return nil
}

func NewEventService(repo repository.EventRepository) EventService {
	return EventService{
		Repository: repo,
	}
}
