package controller

import (
	"calendar.com/pkg/domain/service"
)

type Controller struct {
	EventService service.Event
	AuthService  service.Authorization
}

func NewController(es service.Event, as service.Authorization) *Controller {
	return &Controller{
		EventService: es,
		AuthService:  as,
	}
}
