package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"calendar.com/pkg/logger"
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
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
}

func (c *Client) EventUpdate(w http.ResponseWriter, r *http.Request) {
	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{Name: "id"}.Error(), "event-update")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorUnhandledPathParameter{Name: "id"})
	}

	err := c.EventOperations.Update(r.Context(), w, r.Body, eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-update")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
}

func (c *Client) EventDelete(w http.ResponseWriter, r *http.Request) {
	eventId, ok := mux.Vars(r)["id"]
	if !ok || eventId == "" {
		logger.NewLogger().Write(logger.Error, ErrorUnhandledPathParameter{Name: "id"}.Error(), "event-delete")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: ErrorUnhandledPathParameter{Name: "id"}.Error()})
	}

	err := c.EventOperations.Delete(r.Context(), w, eventId)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-delete")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
}
