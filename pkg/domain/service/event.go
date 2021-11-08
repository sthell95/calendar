package service

import (
	"calendar.com/pkg/storage"
)

type Event interface {
}

type EventService struct{}

func NewEventService(repo *storage.Repository) EventService {
	return EventService{}
}
