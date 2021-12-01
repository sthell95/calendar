package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/storage"
)

type eventPut struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string        `gorm:"type:varchar(100);not null"`
	Description string        `gorm:"type:text"`
	Timezone    string        `gorm:"type:varchar;default 'Europe/Riga'"`
	Time        *time.Time    `gorm:"type:timestamp; not null"`
	Duration    time.Duration `gorm:"type:time not null"`
	User        uuid.UUID     `gorm:"type:uuid;not null"`
	Notes       []entity.Note `gorm:"foreignKey:EventID"`
}

type eventGet struct {
	ID          string        `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string        `gorm:"type:varchar(100);not null"`
	Description string        `gorm:"type:text"`
	Timezone    string        `gorm:"type:varchar;default 'Europe/Riga'"`
	Time        *time.Time    `gorm:"type:timestamp; not null"`
	Duration    time.Duration `gorm:"type:time not null"`
	User        uuid.UUID     `gorm:"type:uuid;not null"`
	Notes       []entity.Note `gorm:"foreignKey:EventID"`
}

type eventDelete struct {
	ID   string    `gorm:"type:uuid;default:uuid_generate_v4()"`
	User uuid.UUID `gorm:"type:uuid;not null"`
}

type EventRepository interface {
	Create(context.Context, *entity.Event) error
	Update(context.Context, *entity.Event) error
	Delete(context.Context, *entity.Event) error
	FindOneById(*uuid.UUID) (*entity.Event, error)
}

type EventModel struct{}

type Event struct {
	repo storage.Repository
}

func (ev *Event) Create(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "create-event-repository")
	defer span.Finish()

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

	e.ID = m.ID
	return nil
}

func (ev *Event) Update(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update-event-repository")
	defer span.Finish()

	m := &eventPut{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		Notes:       e.Notes,
		User:        e.User.ID,
	}

	condition := fmt.Sprintf(`"user" = '%v'`, e.User.ID)
	err := ev.repo.Update(m, EventModel{}, condition)
	if err != nil {
		return err
	}

	return nil
}

func (ev *eventPut) BeforeUpdate(tx *gorm.DB) error {
	return tx.Where(fmt.Sprintf("event_id = '%v'", ev.ID)).Delete(&entity.Note{}).Error
}

func (e *Event) FindOneById(id *uuid.UUID) (*entity.Event, error) {
	var model *eventGet
	err := e.repo.FindById(model, *id)
	if err != nil {
		return nil, err
	}

	event := entity.Event{
		Title:       model.Title,
		Description: "",
		Timezone:    "",
		Time:        nil,
		Duration:    0,
		User:        entity.User{},
		Notes:       nil,
	}

	return &event, nil
}

func (e *Event) Delete(ctx context.Context, event *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete-event-repository")
	defer span.Finish()

	c := fmt.Sprintf(`"user" = '%v'`, event.User.ID)
	m := eventDelete{
		ID:   event.ID.String(),
		User: event.ID,
	}
	return e.repo.Delete(m, EventModel{}, c)
}

func NewEventRepository(repo storage.Repository) *Event {
	return &Event{
		repo: repo,
	}
}

func (EventModel) GetTable() string {
	return "events"
}
