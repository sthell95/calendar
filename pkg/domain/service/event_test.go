package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
)

func TestEventService_Create(t *testing.T) {
	tests := []struct {
		name    string
		event   func() entity.Event
		mock    func(t *testing.T, e entity.Event) repository.EventRepository
		wantErr error
	}{
		{
			name: "Valid",
			event: func() entity.Event {
				add := time.Now().Add(time.Hour * 1)

				return entity.Event{Time: &add}
			},
			mock: func(t *testing.T, e entity.Event) repository.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := repository.NewMockEventRepository(ctrl)
				eventRepository.
					EXPECT().
					Create(&e).
					Return(nil)

				return eventRepository
			},
			wantErr: nil,
		},
		{
			name: "Invalid",
			event: func() entity.Event {
				add := time.Now().Add(-time.Hour * 1)

				return entity.Event{Time: &add}
			},
			mock: func(t *testing.T, e entity.Event) repository.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := repository.NewMockEventRepository(ctrl)

				return eventRepository
			},
			wantErr: IncorrectTime{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.event()
			eventRepository := tt.mock(t, event)
			es := &EventService{Repository: eventRepository}

			err := es.Create(&event)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_Update(t *testing.T) {
	tests := []struct {
		name    string
		event   func() entity.Event
		mock    func(t *testing.T, e entity.Event) repository.EventRepository
		wantErr error
	}{
		{
			name: "Valid",
			event: func() entity.Event {
				add := time.Now().Add(time.Hour * 1)

				return entity.Event{Time: &add}
			},
			mock: func(t *testing.T, e entity.Event) repository.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := repository.NewMockEventRepository(ctrl)
				eventRepository.
					EXPECT().
					Update(&e).
					Return(nil)

				return eventRepository
			},
			wantErr: nil,
		},
		{
			name: "Invalid",
			event: func() entity.Event {
				add := time.Now().Add(-time.Hour * 1)

				return entity.Event{Time: &add}
			},
			mock: func(t *testing.T, e entity.Event) repository.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := repository.NewMockEventRepository(ctrl)

				return eventRepository
			},
			wantErr: IncorrectTime{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.event()
			eventRepository := tt.mock(t, event)
			es := &EventService{Repository: eventRepository}

			err := es.Update(&event)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
