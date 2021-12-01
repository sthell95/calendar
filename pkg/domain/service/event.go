package service

import (
	"calendar.com/pkg/domain/repository/postgres"
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/gofrs/uuid"

	"calendar.com/pkg/domain/entity"
)

type Event interface {
	Create(context.Context, *entity.Event) error
	Update(context.Context, *entity.Event) error
	Delete(context.Context, *entity.Event) error
}

type Validators interface {
	IsAuthor(userId *uuid.UUID, eventUserId *uuid.UUID) bool
	ValidateTime(*time.Time) bool
}

type EventService struct {
	Repository postgres.EventRepository
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

func (es *EventService) Create(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "create-event-process")
	defer span.Finish()

	validator := ValidateEntity{}
	if ok := validator.ValidateTime(e.Time); !ok {
		return IncorrectTime{}
	}

	return es.Repository.Create(ctx, e)
}

func (es *EventService) Update(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update-event-process")
	defer span.Finish()

	validator := ValidateEntity{}
	if ok := validator.ValidateTime(e.Time); !ok {
		return IncorrectTime{}
	}

	return es.Repository.Update(ctx, e)
}

func (es *EventService) Delete(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete-event-process")
	defer span.Finish()

	return es.Repository.Delete(ctx, e)
}

func (*ValidateEntity) IsAuthor(userId *uuid.UUID, eventUserId *uuid.UUID) bool {
	return userId.String() == eventUserId.String()
}

func (*ValidateEntity) ValidateTime(t *time.Time) bool {
	return t.After(time.Now())
}

func NewEventService(repo postgres.EventRepository) *EventService {
	return &EventService{
		Repository: repo,
	}
}
