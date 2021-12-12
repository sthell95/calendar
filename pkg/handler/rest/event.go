package rest

import (
	"calendar.com/middleware"
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"calendar.com/pkg/logger"
)

type EventOperations interface {
	Create(ctx context.Context, e *entity.Event) error
	Update(ctx context.Context, w io.Writer, r io.Reader, eventId string) error
	Delete(ctx context.Context, w io.Writer, eventId string) error
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

type Client struct {
	AuthOperations
	EventOperations
}

type Error struct {
	Message string `json:"message"`
}

type ErrorUserContext struct{}

type ErrorUnhandledPathParameter struct {
	Name  string
	Value string
}

func NewClient(ao AuthOperations, eo EventOperations) *Client {
	return &Client{
		ao,
		eo,
	}
}

func (*ErrorUserContext) Error() string {
	return "User does not exists in the context"
}

func (e ErrorUnhandledPathParameter) Error() string {
	return fmt.Sprintf("Not found parameter: %v, value: %v in request path", e.Name, e.Value)
}

func (c *Client) EventCreate(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "create-event")
	defer span.Finish()

	var event RequestEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-create")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	entityEvent, err := event.RequestToEntity(ctx)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-create")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	err = c.EventOperations.Create(r.Context(), entityEvent)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-create")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	re := &ResponseEvent{}
	re.EntityToResponse(ctx, *entityEvent)
	response.NewPrint().PrettyPrint(w, re)
}

func (c *Client) EventUpdate(w http.ResponseWriter, r *http.Request) {
	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{Name: "id"}.Error(), "event-update")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorUnhandledPathParameter{Name: "id"})
	}

	err := c.EventOperations.Update(r.Context(), w, r.Body, eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-update")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
}

func (c *Client) EventDelete(w http.ResponseWriter, r *http.Request) {
	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{Name: "id"}.Error(), "event-delete")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: ErrorUnhandledPathParameter{Name: "id"}.Error()})
	}

	err := c.EventOperations.Delete(r.Context(), w, eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-delete")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
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
