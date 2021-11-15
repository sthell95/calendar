package service

import (
	"github.com/gofrs/uuid"

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
	eventPut, err := es.Repository.Create(e)
	e.ID = uuid.FromStringOrNil(eventPut.ID)
	return err
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{
		Repository: repo,
	}
}
