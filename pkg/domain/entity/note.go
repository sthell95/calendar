package entity

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"calendar.com/pkg/logger"
)

type Note struct {
	ID   string `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Note string `json:"note" gorm:"type:varchar(60); not null"`
}

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event")
		return err
	}
	tx.Set("ID", id.String())
	return nil
}
