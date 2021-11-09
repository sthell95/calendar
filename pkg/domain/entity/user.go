package entity

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Login    string    `gorm:"type:varchar"`
	Password string    `gorm:"type:varchar"`
	Timezone string    `gorm:"type:varchar"`
}
