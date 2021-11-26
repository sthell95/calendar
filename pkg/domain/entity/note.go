package entity

import "github.com/gofrs/uuid"

type Note struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Note    string    `json:"note" gorm:"type:varchar(60); not null"`
	EventID uuid.UUID `json:"event_id"`
}
