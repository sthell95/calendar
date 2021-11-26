package service

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
)

type Event interface {
	Create(*entity.Event) error
	Update(*entity.Event) error
	Delete(*uuid.UUID) error
}

type Validators interface {
	IsAuthor(userId *uuid.UUID, eventUserId *uuid.UUID) bool
	ValidateTime(*time.Time) bool
}

type EventService struct {
	Repository repository.EventRepository
}

type ValidateEntity struct{}

type IncorrectTime struct{}

func (IncorrectTime) Error() string {
	return fmt.Sprintf("Event time is not correct please choose time in the future")
}

type Forbidden struct{}

func (Forbidden) Error() string {
	return fmt.Sprintf("You don't have permissions for this event")
}

type EventNotFound struct{}

func (EventNotFound) Error() string {
	return fmt.Sprintf("Event not found")
}

func (es *EventService) Create(e *entity.Event) error {
	validator := ValidateEntity{}
	if ok := validator.ValidateTime(e.Time); !ok {
		return IncorrectTime{}
	}

	return es.Repository.Create(e)
}

func (es *EventService) Update(e *entity.Event) error {
	validator := ValidateEntity{}
	if ok := validator.ValidateTime(e.Time); !ok {
		return IncorrectTime{}
	}

	return es.Repository.Update(e)
}

func (es *EventService) Delete(eventId *uuid.UUID) error {
	event, err := es.Repository.FindOneById(eventId)
	if err != nil {
		return err
	}

	return es.Repository.Delete(event)
}

func (*ValidateEntity) IsAuthor(userId *uuid.UUID, eventUserId *uuid.UUID) bool {
	return userId.String() == eventUserId.String()
}

func (*ValidateEntity) ValidateTime(t *time.Time) bool {
	return t.After(time.Now())
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{
		Repository: repo,
	}
}
