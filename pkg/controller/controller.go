package controller

import (
	"calendar.com/pkg/domain/repository"
	"calendar.com/pkg/domain/service"
)

type Controller struct {
	Services   *service.Services
	Repository *repository.Repository
}

func NewController(services *service.Services, repos *repository.Repository) *Controller {
	return &Controller{
		Services:   services,
		Repository: repos,
	}
}
