package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

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
				mock.EXPECT().Create(event).Return(nil)

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
				mock.EXPECT().Create(event).Return(errors.New("Could not create event"))

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
	type fields struct {
		ID          string
		Title       string
		Description string
		Timezone    string
		Time        string
		Duration    int32
		Notes       []string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &RequestEvent{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Timezone:    tt.fields.Timezone,
				Time:        tt.fields.Time,
				Duration:    tt.fields.Duration,
				Notes:       tt.fields.Notes,
			}
			got, err := re.RequestToEntity(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestToEntity() got = %v, want %v", got, tt.want)
			}
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
	type args struct {
		e entity.Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ResponseEvent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &ResponseEvent{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Timezone:    tt.fields.Timezone,
				Time:        tt.fields.Time,
				Duration:    tt.fields.Duration,
				Notes:       tt.fields.Notes,
			}
			re.EntityToResponse(tt.args.e)
			//if got := re.EntityToResponse(tt.args.e); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("EntityToResponse() = %v, want %v", got, tt.want)
			//}
		})
	}
}
