package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"

	"calendar.com/middleware"
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

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

func (re *RequestEvent) RequestToEntity(ctx context.Context) (*entity.Event, error) {
	t, err := time.Parse(entity.ISOLayout, re.Time)
	if err != nil {
		return nil, err
	}

	d, err := time.ParseDuration(fmt.Sprintf("%vs", re.Duration))
	if err != nil {
		return nil, err
	}

	if userId, ok := ctx.Value(middleware.UserId).(uuid.UUID); ok {
		e := &entity.Event{
			Title:       re.Title,
			Description: re.Description,
			Timezone:    re.Timezone,
			Time:        &t,
			Duration:    d,
			User: entity.User{
				ID: userId,
			},
		}
		for _, v := range re.Notes {
			e.Notes = append(e.Notes, entity.Note{Note: v})
		}

		return e, nil
	}

	return nil, errors.New("User does not exists in the context")
}

func (re *ResponseEvent) EntityToResponse(e entity.Event) *ResponseEvent {
	re = &ResponseEvent{
		ID:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time.Format(entity.ISOLayout),
		Duration:    int32(e.Duration.Seconds()),
	}
	for _, note := range e.Notes {
		re.Notes = append(re.Notes, note.Note)
	}

	return re
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

	re := ResponseEvent{}
	resp := re.EntityToResponse(*entityEvent)
	response.NewPrint().PrettyPrint(w, resp)
}
