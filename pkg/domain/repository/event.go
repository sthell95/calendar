package repository

import (
	"time"

	"github.com/gofrs/uuid"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/storage"
)

type eventPut struct {
	ID          string        `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string        `gorm:"type:varchar(100);not null"`
	Description string        `gorm:"type:text"`
	Timezone    string        `gorm:"type:varchar;default 'Europe/Riga'"`
	Time        *time.Time    `gorm:"type:timestamp; not null"`
	Duration    time.Duration `gorm:"type:time not null"`
	User        uuid.UUID     `gorm:"type:uuid;not null"`
	Notes       []entity.Note `gorm:"foreignKey:EventID"`
}

type EventRepository interface {
	Create(*entity.Event) error
	Update(*entity.Event) error
	FindOneById(string) (*entity.Event, error)
}

type EventModel struct{}

type Event struct {
	repo storage.Repository
}

func (ev *Event) Create(e *entity.Event) error {
	t := time.Now()
	m := &eventPut{
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        &t,
		Duration:    e.Duration,
		Notes:       e.Notes,
		User:        e.User.ID,
	}

	err := ev.repo.Create(m, EventModel{})
	if err != nil {
		return err
	}

	e.ID = uuid.FromStringOrNil(m.ID)
	return nil
}

func (ev *Event) Update(e *entity.Event) error {
	m := &eventPut{
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		Notes:       e.Notes,
		User:        e.User.ID,
	}

	err := ev.repo.Update(m, EventModel{})
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) FindOneById(id string) (*entity.Event, error) {
	return nil, nil
}

func NewEventRepository(repo storage.Repository) *Event {
	return &Event{
		repo: repo,
	}
}

func (EventModel) GetTable() string {
	return "events"
}
