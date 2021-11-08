package controller

import (
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/logger"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
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

	timestamp := strconv.ParseInt()
	t := time.Unix(event.Time, 0)
}
