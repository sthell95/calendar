package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	var event entity.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	err = c.Repository.Event.Create(&event)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
		return
	}

	response.NewPrint().PrettyPrint(w, event)
}
