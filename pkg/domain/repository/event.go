package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"calendar.com/pkg/domain/entity"
)

type EventRepository interface {
	Create(context.Context, *entity.Event) error
	Update(context.Context, *entity.Event) error
	Delete(context.Context, *entity.Event) error
	FindOneById(*uuid.UUID) (*entity.Event, error)
}

type EventRepo struct {
	Repos []EventRepository
}

func NewEventRepository(repos ...EventRepository) *EventRepo {
	return &EventRepo{Repos: repos}
}

func (r *EventRepo) Create(ctx context.Context, event *entity.Event) error {
	var err error
	event.ID = uuid.New()

	for i := range r.Repos {
		if err = r.Repos[i].Create(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepo) Update(ctx context.Context, event *entity.Event) error {
	for i := range r.Repos {
		if err := r.Repos[i].Update(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepo) Delete(ctx context.Context, event *entity.Event) error {
	for i := range r.Repos {
		if err := r.Repos[i].Delete(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepo) FindOneById(id *uuid.UUID) (*entity.Event, error) {
	for i := range r.Repos {
		return r.Repos[i].FindOneById(id)
	}

	return nil, errors.New("no providers")
}
