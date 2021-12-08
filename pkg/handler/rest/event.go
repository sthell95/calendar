package rest

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type EventOperations interface {
	Create(ctx context.Context, w io.Writer, r io.Reader) error
	Update(ctx context.Context, w io.Writer, r io.Reader, eventId string) error
	Delete(ctx context.Context, w io.Writer, eventId string) error
}

type Client struct {
	AuthOperations
	EventOperations
}

type Error struct {
	Message string `json:"message"`
}

type ErrorUnhandledPathParameter struct {
	Name  string
	Value string
}

func NewClient(ao AuthOperations, eo EventOperations) *Client {
	return &Client{
		ao,
		eo,
	}
}

func (e ErrorUnhandledPathParameter) Error() string {
	return fmt.Sprintf("Not found parameter: %v, value: %v in request path", e.Name, e.Value)
}

func (c *Client) EventCreate(w http.ResponseWriter, r *http.Request) {
	err := c.EventOperations.Create(r.Context(), w, r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-create")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}

func (c *Client) EventUpdate(w http.ResponseWriter, r *http.Request) {
	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{Name: "id"}.Error(), "event-update")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: ErrorUnhandledPathParameter{Name: "id"}.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}

	err := c.EventOperations.Update(r.Context(), w, r.Body, eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-update")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}

func (c *Client) EventDelete(w http.ResponseWriter, r *http.Request) {
	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{Name: "id"}.Error(), "event-delete")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: ErrorUnhandledPathParameter{Name: "id"}.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}

	err := c.EventOperations.Delete(r.Context(), w, eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-delete")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}
