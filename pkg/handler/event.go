package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"calendar.com/pkg/domain/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"

	"calendar.com/middleware"
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type EventHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type RequestEvent struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Timezone    string   `json:"timezone"`
	Time        string   `json:"time"`
	Duration    int32    `json:"duration"`
	Notes       []string `json:"notes"`
}

type ResponseEvent struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Timezone    string   `json:"timezone"`
	Time        string   `json:"time"`
	Duration    int32    `json:"duration"`
	Notes       []string `json:"notes"`
}

type Event struct {
	service.Event
}

type ErrorUserContext struct{}

type ErrorUnhandledPathParameter struct {
	Name  string
	Value string
}

func NewEventOperations(s service.Event) *Event {
	return &Event{s}
}

func (*ErrorUserContext) Error() string {
	return "User does not exists in the context"
}

func (e ErrorUnhandledPathParameter) Error() string {
	return fmt.Sprintf("Not found parameter: %v, value: %v in request path", e.Name, e.Value)
}

func (re *RequestEvent) RequestToEntity(ctx context.Context) (*entity.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "request-to-entity")
	defer span.Finish()

	t, err := time.Parse(entity.ISOLayout, re.Time)
	if err != nil {
		return nil, err
	}

	d, err := time.ParseDuration(fmt.Sprintf("%vs", re.Duration))
	if err != nil {
		return nil, err
	}

	var e entity.Event
	if eventId, ok := ctx.Value(entity.EventIdKey).(uuid.UUID); ok {
		e.ID = eventId
	}

	if userId, ok := ctx.Value(middleware.UserId).(uuid.UUID); ok {
		e.Title = re.Title
		e.Description = re.Description
		e.Timezone = re.Timezone
		e.Time = &t
		e.Duration = d
		e.User = entity.User{ID: userId}

		e.Notes = make([]entity.Note, len(re.Notes), len(re.Notes))
		for i, v := range re.Notes {
			e.Notes[i] = entity.Note{Note: v, EventID: e.ID}
		}

		return &e, nil
	}

	return nil, &ErrorUserContext{}
}

func (re *ResponseEvent) EntityToResponse(ctx context.Context, e entity.Event) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "entity-to-response")
	defer span.Finish()

	re.ID = e.ID.String()
	re.Title = e.Title
	re.Description = e.Description
	re.Timezone = e.Timezone
	re.Time = e.Time.Format(entity.ISOLayout)
	re.Duration = int32(e.Duration.Seconds())

	for _, note := range e.Notes {
		re.Notes = append(re.Notes, note.Note)
	}
}

func (c *Event) Create(w http.ResponseWriter, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "create-event")
	defer span.Finish()

	var event RequestEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return err
	}

	entityEvent, err := event.RequestToEntity(r.Context())
	if err != nil {
		return err
	}

	err = c.Event.Create(ctx, entityEvent)
	if err != nil {
		return err
	}

	re := &ResponseEvent{}
	re.EntityToResponse(ctx, *entityEvent)
	response.NewPrint().PrettyPrint(w, re)
	return nil
}

func (c *Event) Update(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "update-event")
	defer span.Finish()

	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{
			Name:  "id",
			Value: eventId,
		}.Error(), "Update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: ErrorUnhandledPathParameter{
			Name:  "id",
			Value: eventId,
		}.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	var event RequestEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "Update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	id, err := uuid.Parse(eventId)
	ctx = context.WithValue(ctx, entity.EventIdKey, id)
	entityEvent, err := event.RequestToEntity(ctx)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "Update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	err = c.Event.Update(ctx, entityEvent)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "Update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	re := &ResponseEvent{}
	re.EntityToResponse(ctx, *entityEvent)
	response.NewPrint().PrettyPrint(w, re)
}

func (c *Event) Delete(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "delete-event")
	span.SetTag("this", 999999999999)
	defer span.Finish()

	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{}.Error(), "delete-event")
		response.NewPrint().PrettyPrint(w, Error{Message: ErrorUnhandledPathParameter{}.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	id, err := uuid.Parse(eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "delete-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	userId := r.Context().Value(middleware.UserId).(uuid.UUID)
	e := entity.Event{ID: id, User: entity.User{ID: userId}}
	err = c.Event.Delete(ctx, &e)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "delete-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
