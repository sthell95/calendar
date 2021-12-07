package handler

import (
	"calendar.com/pkg/domain/service"
)

type Controller struct {
	EventService service.Event
	AuthService  service.Authorization
}

type Error struct {
	Message string `json:"message"`
}

func NewController(es service.Event, as service.Authorization) *Controller {
	return &Controller{
		EventService: es,
		AuthService:  as,
	}
}
