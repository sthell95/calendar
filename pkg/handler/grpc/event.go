package grpc

import (
	"calendar.com/pkg/domain/entity"
	pg "calendar.com/proto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

type EventOperations interface {
	Create(ctx context.Context, e *entity.Event) error
	Update(ctx context.Context, w io.Writer, r io.Reader, eventId string) error
	Delete(ctx context.Context, w io.Writer, eventId string) error
}

type EventHandler struct {
	pg.UnimplementedEventServiceServer
	EventOperations EventOperations
}

func (h *EventHandler) Create(ctx context.Context, event *pg.Event) (*pg.Event, error) {
	model, err := eventGrpcToDomain(ctx, event)
	if err != nil {
		return nil, err
	}

	err = h.EventOperations.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *EventHandler) Update(ctx context.Context, event *pg.Event) (*pg.Event, error) {
	return nil, nil
}

func (h *EventHandler) Delete(ctx context.Context, eventId *pg.EventId) (*pg.DeleteResponse, error) {
	return nil, nil
}

func eventGrpcToDomain(ctx context.Context, re *pg.Event) (*entity.Event, error) {
	id, err := uuid.Parse(re.Id)
	if err != nil {
		return nil, err
	}

	eventTime, err := time.Parse(entity.ISOLayout, re.Time)
	if err != nil {
		return nil, err
	}

	eventDuration, err := time.ParseDuration(fmt.Sprint(re.Duration))
	if err != nil {
		return nil, err
	}

	if eventId, ok := ctx.Value(entity.EventIdKey).(uuid.UUID); ok {
		re.Id = eventId.String()
	}

	return &entity.Event{
		ID:          id,
		Title:       re.Title,
		Description: re.Description,
		Timezone:    re.Timezone,
		Time:        &eventTime,
		Duration:    eventDuration,
		User:        entity.User{},
		Notes:       nil,
	}, nil
}

func DomainToEventGrpc(e *entity.Event) *pg.Event {
	event := &pg.Event{
		Id:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time.Format(entity.ISOLayout),
		Duration:    int64(e.Duration.Seconds()),
	}

	var n pg.Note
	for _, note := range e.Notes {
		n.Note = note.Note
		event.Notes = append(event.Notes, &n)
	}

	return event
}
