package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"calendar.com/pkg/response"

	"calendar.com/pkg/domain/entity"

	"calendar.com/pkg/logger"
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

	d, err := time.ParseDuration(fmt.Sprintf("%vs", event.Duration))
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
	}
	for _, v := range event.Notes {
		entityEvent.Notes = append(entityEvent.Notes, entity.Note{Note: v})
	}
	err = c.EventService.Create(&entityEvent)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	resp := &ResponseEvent{
		ID:          entityEvent.ID.String(),
		Title:       entityEvent.Title,
		Description: entityEvent.Description,
		Timezone:    entityEvent.Timezone,
		Time:        entityEvent.Time.Format(entity.ISOLayout),
		Duration:    int32(entityEvent.Duration.Seconds()),
	}
	for _, note := range entityEvent.Notes {
		resp.Notes = append(resp.Notes, note.Note)
	}
	response.NewPrint().PrettyPrint(w, resp)
}
