package repository

import (
	"time"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/storage"
)

type EventPut struct {
	ID          string        `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string        `gorm:"type:varchar(100);not null"`
	Description string        `gorm:"type:text"`
	Timezone    string        `gorm:"type:varchar;default 'Europe/Riga'"`
	Time        *time.Time    `gorm:"type:timestamp; not null"`
	Duration    time.Duration `gorm:"type:time not null"`
	//User        User       `json:"-"`
	Notes []string `gorm:"type:varchar []"`
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

	_ = ev.repo.Create(&EventPut{
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		Notes:       e.Notes,
	})

	t := time.Now()
	m := EventPut{
		ID:          "d3f09345-e565-416a-ad3a-0f4fe51f8842",
		Title:       "Some title",
		Description: "Desfription",
		Timezone:    "Africa",
		Time:        &t,
		Duration:    time.Duration(3600),
		Notes: []string{
			"First",
			"Second",
			"Last",
		},
	}

	return ev.repo.Create(&m)
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
