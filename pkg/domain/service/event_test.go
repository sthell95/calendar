package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	postgres2 "calendar.com/pkg/storage/postgresdb"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"

	"calendar.com/pkg/domain/entity"
)

func TestEventService_Create(t *testing.T) {
	tests := []struct {
		name    string
		event   func() entity.Event
		mock    func(t *testing.T, e entity.Event) postgres2.EventRepository
		wantErr error
	}{
		{
			name: "Valid",
			event: func() entity.Event {
				add := time.Now().Add(time.Hour * 1)

				return entity.Event{Time: &add}
			},
			mock: func(t *testing.T, e entity.Event) postgres2.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := postgres2.NewMockEventRepository(ctrl)
				eventRepository.
					EXPECT().
					Create(gomock.Any(), &e).Do(func(a, b interface{}) {
					fmt.Println(a, b)
				}).
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
			mock: func(t *testing.T, e entity.Event) postgres2.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := postgres2.NewMockEventRepository(ctrl)

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

			err := es.Create(context.Background(), &event)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_Update(t *testing.T) {
	tests := []struct {
		name    string
		event   func() entity.Event
		mock    func(t *testing.T, e entity.Event) postgres2.EventRepository
		wantErr error
	}{
		{
			name: "Valid",
			event: func() entity.Event {
				add := time.Now().Add(time.Hour * 1)

				return entity.Event{Time: &add}
			},
			mock: func(t *testing.T, e entity.Event) postgres2.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := postgres2.NewMockEventRepository(ctrl)
				eventRepository.
					EXPECT().
					Update(gomock.Any(), &e).
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
			mock: func(t *testing.T, e entity.Event) postgres2.EventRepository {
				ctrl := gomock.NewController(t)
				eventRepository := postgres2.NewMockEventRepository(ctrl)

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

			err := es.Update(context.Background(), &event)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
