package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"calendar.com/middleware"
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/service"
)

func TestController_Create(t *testing.T) {
	tests := []struct {
		name         string
		mock         func(*testing.T, *RequestEvent, context.Context) service.Event
		requestEvent func() *RequestEvent
		event        *entity.Event
		ctx          func() context.Context
		want         string
	}{
		{
			name: "Valid",
			mock: func(t *testing.T, e *RequestEvent, ctx context.Context) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)
				event, _ := e.RequestToEntity(ctx)
				mock.EXPECT().Create(gomock.Any(), event).Return(nil)

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			requestEvent: func() *RequestEvent {
				return &RequestEvent{
					Title:       "Birthday",
					Description: "Need to buy a gift",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    0,
					Notes:       nil,
				}
			},
			want: `{"id":"00000000-0000-0000-0000-000000000000","title":"Birthday","description":"Need to buy a gift","timezone":"America/Chicago","time":"2021-12-10T15:04:05.000Z","duration":0,"notes":null}`,
		},
		{
			name: "User does not exist in context",
			mock: func(t *testing.T, e *RequestEvent, ctx context.Context) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)

				return mock
			},
			ctx: func() context.Context {
				return context.Background()
			},
			requestEvent: func() *RequestEvent {
				return &RequestEvent{
					Title:       "Birthday",
					Description: "Need to buy a gift",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    0,
					Notes:       nil,
				}
			},
			want: `{"message":"User does not exists in the context"}`,
		},
		{
			name: "User does not exist in context",
			mock: func(t *testing.T, e *RequestEvent, ctx context.Context) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)
				event, _ := e.RequestToEntity(ctx)
				mock.EXPECT().Create(gomock.Any(), event).Return(errors.New("Could not create event"))

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			requestEvent: func() *RequestEvent {
				return &RequestEvent{
					Title:       "Birthday",
					Description: "Need to buy a gift",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    0,
					Notes:       nil,
				}
			},
			want: `{"message":"Could not create event"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.requestEvent()
			ctx := tt.ctx()
			eventService := tt.mock(t, event, ctx)
			c := &Controller{EventService: eventService}
			w := httptest.NewRecorder()
			body, _ := json.Marshal(event)

			r := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(body))
			r = r.WithContext(ctx)
			c.Create(w, r)
			responseBody, _ := io.ReadAll(w.Body)
			eventResponse := strings.Trim(string(responseBody), "\n")

			require.Equal(t, tt.want, eventResponse)
		})
	}
}

func TestRequestEvent_RequestToEntity(t *testing.T) {
	tests := []struct {
		name    string
		fields  func(time.Time) *RequestEvent
		ctx     func() *context.Context
		want    func() *entity.Event
		wantErr error
	}{
		{
			name: "Valid",
			fields: func(t time.Time) *RequestEvent {
				return &RequestEvent{
					Title:       "Test title",
					Description: "description",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    3600,
					Notes: []string{
						"First", "Second", "Last",
					},
				}
			},
			ctx: func() *context.Context {
				id, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, id)
				return &ctx
			},
			want: func() *entity.Event {
				t, _ := time.Parse(entity.ISOLayout, "2021-12-10T15:04:05.000Z")
				d, _ := time.ParseDuration(fmt.Sprintf("%vs", 3600))
				id, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				return &entity.Event{
					Title:       "Test title",
					Description: "description",
					Timezone:    "America/Chicago",
					Time:        &t,
					Duration:    d,
					User: entity.User{
						ID: id,
					},
					Notes: []entity.Note{
						{Note: "First"},
						{Note: "Second"},
						{Note: "Last"},
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "Parse time error",
			fields: func(t time.Time) *RequestEvent {
				return &RequestEvent{Time: "2021"}
			},
			ctx: func() *context.Context {
				ctx := context.Background()
				return &ctx
			},
			want: func() *entity.Event {
				return nil
			},
			wantErr: &time.ParseError{},
		},
		{
			name: "Parse duration error",
			fields: func(t time.Time) *RequestEvent {
				return &RequestEvent{
					Title:       "Test title",
					Description: "description",
					Timezone:    "America/Chicago",
					Duration:    -1,
					Notes:       []string{"First", "Second", "Last"},
				}
			},
			ctx: func() *context.Context {
				id, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, id)
				return &ctx
			},
			want: func() *entity.Event {
				return nil
			},
			wantErr: &time.ParseError{},
		},
		{
			name: "Invalid context",
			fields: func(t time.Time) *RequestEvent {
				return &RequestEvent{
					Title:       "Test title",
					Description: "description",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    3600,
					Notes: []string{
						"First", "Second", "Last",
					},
				}
			},
			ctx: func() *context.Context {
				ctx := context.Background()
				return &ctx
			},
			want: func() *entity.Event {
				return nil
			},
			wantErr: &ErrorUserContext{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventTime := time.Now().Add(time.Hour)
			re := tt.fields(eventTime)
			got, err := re.RequestToEntity(*tt.ctx())
			if err != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			}

			isEqual := reflect.DeepEqual(got, tt.want())
			require.Equal(t, isEqual, true)
		})
	}
}

func TestResponseEvent_EntityToResponse(t *testing.T) {
	type fields struct {
		ID          string
		Title       string
		Description string
		Timezone    string
		Time        string
		Duration    int32
		Notes       []string
	}
	tests := []struct {
		name        string
		fields      fields
		eventEntity func() *entity.Event
		want        *ResponseEvent
	}{
		{
			name: "Valid",
			eventEntity: func() *entity.Event {
				id, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				t, _ := time.Parse(entity.ISOLayout, "2021-12-10T15:04:05.000Z")
				d, _ := time.ParseDuration(fmt.Sprintf("%vs", 3600))
				return &entity.Event{
					ID:          id,
					Title:       "Test",
					Description: "test description",
					Timezone:    "America/Chicago",
					Time:        &t,
					Duration:    d,
					User:        entity.User{},
					Notes: []entity.Note{
						{Note: "First"},
						{Note: "Second"},
						{Note: "Last"},
					},
				}
			},
			want: &ResponseEvent{
				ID:          "62b45338-ea71-4eaa-b5dd-0b29c752ad1c",
				Title:       "Test",
				Description: "test description",
				Timezone:    "America/Chicago",
				Time:        "2021-12-10T15:04:05.000Z",
				Duration:    3600,
				Notes:       []string{"First", "Second", "Last"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ResponseEvent{}
			re.EntityToResponse(context.Background(), *tt.eventEntity())

			require.Equal(t, tt.want, re)
		})
	}
}

func TestController_Update(t *testing.T) {
	tests := []struct {
		name         string
		mock         func(*testing.T, *RequestEvent, context.Context) service.Event
		requestEvent func() *RequestEvent
		ctx          func() context.Context
		want         string
		eventId      string
	}{
		{
			name: "Valid",
			mock: func(t *testing.T, e *RequestEvent, ctx context.Context) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)
				event, _ := e.RequestToEntity(ctx)
				mock.EXPECT().Update(gomock.Any(), event).Return(nil)

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			requestEvent: func() *RequestEvent {
				return &RequestEvent{
					Title:       "Birthday",
					Description: "Need to buy a gift",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    0,
					Notes:       nil,
				}
			},
			want:    `{"id":"00000000-0000-0000-0000-000000000000","title":"Birthday","description":"Need to buy a gift","timezone":"America/Chicago","time":"2021-12-10T15:04:05.000Z","duration":0,"notes":null}`,
			eventId: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "Invalid id",
			mock: func(t *testing.T, e *RequestEvent, ctx context.Context) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			requestEvent: func() *RequestEvent {
				return &RequestEvent{
					Title:       "Birthday",
					Description: "Need to buy a gift",
					Timezone:    "America/Chicago",
					Time:        "2021-12-10T15:04:05.000Z",
					Duration:    0,
					Notes:       nil,
				}
			},
			want:    `{"message":"Not found parameter: id, value:  in request path"}`,
			eventId: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.requestEvent()
			ctx := tt.ctx()
			eventService := tt.mock(t, event, ctx)
			c := &Controller{EventService: eventService}
			w := httptest.NewRecorder()
			body, _ := json.Marshal(event)
			r := httptest.NewRequest(http.MethodPut, "/api/events/:id", bytes.NewReader(body))
			r = r.WithContext(ctx)

			if tt.eventId != "" {
				r = mux.SetURLVars(r, map[string]string{"id": tt.eventId})
			}

			c.Update(w, r)
			responseBody, _ := io.ReadAll(w.Body)
			eventResponse := strings.Trim(string(responseBody), "\n")

			require.Equal(t, tt.want, eventResponse)
		})
	}
}

func TestController_Delete(t *testing.T) {
	tests := []struct {
		name      string
		mock      func(*testing.T, *entity.Event) service.Event
		ctx       func() context.Context
		wantError int
		eventId   string
	}{
		{
			name: "Valid",
			mock: func(t *testing.T, e *entity.Event) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)
				mock.EXPECT().Delete(gomock.Any(), e).Return(nil)

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			wantError: http.StatusNoContent,
			eventId:   "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "Invalid id",
			mock: func(t *testing.T, e *entity.Event) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			wantError: http.StatusBadRequest,
			eventId:   "",
		},
		{
			name: "Invalid deletion process",
			mock: func(t *testing.T, e *entity.Event) service.Event {
				ctrl := gomock.NewController(t)
				mock := service.NewMockEvent(ctrl)
				mock.EXPECT().Delete(gomock.Any(), e).Return(errors.New("Couldn't delete"))

				return mock
			},
			ctx: func() context.Context {
				userId, _ := uuid.FromString("62b45338-ea71-4eaa-b5dd-0b29c752ad1c")
				ctx := context.WithValue(context.Background(), middleware.UserId, userId)

				return ctx
			},
			wantError: http.StatusBadRequest,
			eventId:   "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx()
			eventId, _ := uuid.FromString(tt.eventId)
			e := entity.Event{
				ID:   eventId,
				User: entity.User{ID: ctx.Value(middleware.UserId).(uuid.UUID)},
			}
			eventService := tt.mock(t, &e)
			c := &Controller{EventService: eventService}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/api/events/:id", nil)
			r = r.WithContext(ctx)

			if tt.eventId != "" {
				r = mux.SetURLVars(r, map[string]string{"id": tt.eventId})
			}

			c.Delete(w, r)

			require.Equal(t, tt.wantError, w.Code)
		})
	}
}
