package entity

type Note struct {
	ID      string `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Note    string `json:"note" gorm:"type:varchar(60); not null"`
	EventID string
}
