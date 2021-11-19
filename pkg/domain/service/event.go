package service

import (
	"fmt"
	"time"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
)

type Event interface {
	Create(event *entity.Event) error
	Update(event *entity.Event) error
}

type EventService struct {
	Repository repository.EventRepository
}

type IncorrectTime struct{}

func (IncorrectTime) Error() string {
	return fmt.Sprintf("Event time is not correct please choose time in the future")
}

func (es *EventService) Create(e *entity.Event) error {
	if ok := validateTime(e.Time); !ok {
		return IncorrectTime{}
	}

	return es.Repository.Create(e)
}

func (es *EventService) Update(e *entity.Event) error {
	if ok := validateTime(e.Time); !ok {
		return IncorrectTime{}
	}

	return es.Repository.Update(e)
}

func validateTime(t *time.Time) bool {
	return t.After(time.Now())
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{
		Repository: repo,
	}
}
