package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func (re *RequestEvent) RequestToEntity() (*entity.Event, error) {
	t, err := time.Parse(entity.ISOLayout, re.Time)
	if err != nil {
		return nil, err
	}

	d, err := time.ParseDuration(fmt.Sprintf("%vs", re.Duration))
	if err != nil {
		return nil, err
	}

	e := &entity.Event{
		Title:       re.Title,
		Description: re.Description,
		Timezone:    re.Timezone,
		Time:        &t,
		Duration:    d,
		User:        entity.User{},
	}
	for _, v := range re.Notes {
		e.Notes = append(e.Notes, entity.Note{Note: v})
	}
	return e, nil
}

func (re *ResponseEvent) EntityToResponse(e entity.Event) {
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
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var event RequestEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		response.NewPrint().PrettyPrint(w, Error{Message: err.Error()}, response.WithCode(http.StatusBadRequest))
		return
	}

	entityEvent, err := event.RequestToEntity()
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

	resp := ResponseEvent{}
	resp.EntityToResponse(*entityEvent)
	response.NewPrint().PrettyPrint(w, resp)
}
