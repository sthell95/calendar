package controller

import (
	"encoding/json"
	"net/http"

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

	err = c.EventService.Create(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	response.NewPrint().PrettyPrint(w, event)
}
