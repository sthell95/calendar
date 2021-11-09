package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"calendar.com/pkg/logger"
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
	Notes       []string
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event")
		return err
	}
	tx.Set("ID", id.String())
	//e.ID = eventId.String()
	return nil
}
