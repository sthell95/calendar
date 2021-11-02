package user

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Login    string    `gorm:"type:varchar" json:"login"`
	Password string    `gorm:"type:varchar" json:"password"`
	Salt     string    `gorm:"type:varchar" json:"salt"`
	Timezone string    `gorm:"type:varchar" json:"timezone"`
}
