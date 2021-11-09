package repository

import (
	"time"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/storage"
)

type eventPut struct {
	ID          string     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string     `gorm:"type:varchar(100);not null"`
	Description string     `gorm:"type:text"`
	Timezone    string     `gorm:"type:varchar;default 'Europe/Riga'"`
	Time        *time.Time `gorm:"type:timestamp; not null"`
	Duration    int32      `gorm:"type:time not null"`
	//User        User       `json:"-"`
	//Notes []Note `gorm:"type:"`
}

type EventRepository interface {
	Create(event *entity.Event) error
	Update(event *entity.Event, id string) (*entity.Event, error)
	FindOneById(string) (*entity.Event, error)
}

type Event struct {
	repo storage.Repository
}

func (ev *Event) Create(e *entity.Event) error {

	return ev.repo.Create(&eventPut{
		ID:          e.ID,
		Title:       e.Title,
		Description: "",
		Timezone:    "",
		Time:        nil,
		Duration:    0,
		//Notes:       nil,
	})
}

func (e *Event) Update(event *entity.Event, id string) (*entity.Event, error) {
	return nil, nil
}

func (e *Event) FindOneById(id string) (*entity.Event, error) {
	return nil, nil
}

func NewEventRepository(repo storage.Repository) *Event {
	return &Event{
		repo: repo,
	}
}
