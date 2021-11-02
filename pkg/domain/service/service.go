package service

import "calendar.com/pkg/domain/repository"

type Services struct {
	Authorization
}

func NewService(repos *repository.Repository) *Services {
	return &Services{}
}
