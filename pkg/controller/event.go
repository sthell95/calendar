package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/gofrs/uuid"

	"calendar.com/middleware"
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type contextKey string

const eventIdKey contextKey = "eventId"

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

type ErrorUserContext struct{}

func (*ErrorUserContext) Error() string {
	return "User does not exists in the context"
}

type ErrorUnhandledPathParameter struct {
	Name  string
	Value string
}

func (e ErrorUnhandledPathParameter) Error() string {
	return fmt.Sprintf("Not found parameter: %v, value: %v in request path", e.Name, e.Value)
}

func (re *RequestEvent) RequestToEntity(ctx context.Context) (*entity.Event, error) {
	t, err := time.Parse(entity.ISOLayout, re.Time)
	if err != nil {
		return nil, err
	}

	d, err := time.ParseDuration(fmt.Sprintf("%vs", re.Duration))
	if err != nil {
		return nil, err
	}
	var e entity.Event

	if eventId, ok := ctx.Value(eventIdKey).(contextKey); ok {
		id := uuid.FromStringOrNil(string(eventId))
		e.ID = id
	}

	if userId, ok := ctx.Value(middleware.UserId).(uuid.UUID); ok {
		e = entity.Event{
			Title:       re.Title,
			Description: re.Description,
			Timezone:    re.Timezone,
			Time:        &t,
			Duration:    d,
			User: entity.User{
				ID: userId,
			},
		}
		e.Notes = make([]entity.Note, len(re.Notes), len(re.Notes))
		for i, v := range re.Notes {
			e.Notes[i] = entity.Note{Note: v}
		}

		return &e, nil
	}

	return nil, &ErrorUserContext{}
}

func (re *ResponseEvent) EntityToResponse(e entity.Event) {
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

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var event RequestEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	entityEvent, err := event.RequestToEntity(r.Context())
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	err = c.EventService.Create(entityEvent)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	re := &ResponseEvent{}
	re.EntityToResponse(*entityEvent)
	response.NewPrint().PrettyPrint(w, re)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{}.Error(), "update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: ErrorUnhandledPathParameter{}.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	var event RequestEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	ctx := context.WithValue(r.Context(), eventIdKey, id)
	entityEvent, err := event.RequestToEntity(ctx)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	err = c.EventService.Update(entityEvent)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "update-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	re := &ResponseEvent{}
	re.EntityToResponse(*entityEvent)
	response.NewPrint().PrettyPrint(w, re)
}
