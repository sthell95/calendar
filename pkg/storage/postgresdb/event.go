package postgresdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"calendar.com/pkg/domain/entity"
)

const eventTable = "events"

type event struct {
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

type eventModel struct{}

func (c *Client) Create(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "create-event-repository")
	defer span.Finish()

	*e.Time = time.Now()
	em := eventModel{}
	m := em.toModel(e)

	err := c.db.Table(eventTable).Create(m).Error
	if err != nil {
		return err
	}

	e.ID = m.ID
	return nil
}

func (c *Client) Update(ctx context.Context, e *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update-event-repository")
	defer span.Finish()

	em := eventModel{}
	m := em.toModel(e)

	condition := fmt.Sprintf(`"user" = '%v'`, e.User.ID)
	err := c.db.Table(eventTable).Where(condition).Updates(m).Error
	if err != nil {
		return err
	}

	return nil
}

func (ev *event) BeforeUpdate(tx *gorm.DB) error {
	return tx.Where(fmt.Sprintf("event_id = '%v'", ev.ID)).Delete(&entity.Note{}).Error
}

func (c *Client) FindOneById(id *uuid.UUID) (*entity.Event, error) {
	var e *event
	err := c.db.First(e, *id).Error
	if err != nil {
		return nil, err
	}

	em := eventModel{}
	return em.toDomainModel(e), nil
}

func (c *Client) Delete(ctx context.Context, event *entity.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete-event-repository")
	defer span.Finish()

	conditions := fmt.Sprintf(`"user" = '%v'`, event.User.ID)
	m := eventDelete{
		ID:   event.ID.String(),
		User: event.ID,
	}
	return c.db.Table(eventTable).Where(conditions).Delete(m).Error
}

func (*eventModel) toModel(e *entity.Event) *event {
	return &event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		User:        e.User.ID,
		Notes:       e.Notes,
	}
}

func (*eventModel) toDomainModel(e *event) *entity.Event {
	return &entity.Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		User:        entity.User{ID: e.User},
		Notes:       e.Notes,
	}
}

func NewRepository(db *gorm.DB) *Client {
	return &Client{
		db: db,
	}
}
