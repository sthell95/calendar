package entity

import (
	"time"

	"github.com/google/uuid"
)

const ISOLayout = "2006-01-02T15:04:05.000Z"

type Event struct {
	ID          uuid.UUID
	Title       string
	Description string
	Timezone    string
	Time        *time.Time
	Duration    time.Duration
	User        User
	Notes       []Note
}

type contextKey string

const EventIdKey contextKey = "eventId"
