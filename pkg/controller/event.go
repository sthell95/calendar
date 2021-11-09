package controller

import (
	"encoding/json"
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

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var event RequestEvent

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	t, err := time.Parse(entity.ISOLayout, event.Time)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	d, err := time.ParseDuration(string(event.Duration))
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}
	entityEvent := entity.Event{
		Title:       event.Title,
		Description: event.Description,
		Timezone:    event.Timezone,
		Time:        &t,
		Duration:    d,
		User:        entity.User{},
		Notes:       event.Notes,
	}
	err = c.EventService.Create(&entityEvent)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	response.NewPrint().PrettyPrint(w, event)
}
