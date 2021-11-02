package controller

import (
	"calendar.com/pkg/domain/service"
)

type Controller struct {
	Services *service.Services
}

func NewController(services *service.Services) *Controller {
	return &Controller{Services: services}
}
