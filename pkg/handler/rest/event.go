package rest

import (
	"context"
	"io"
	"net/http"

	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type EventOperations interface {
	Create(ctx context.Context, w io.Writer, r io.Reader) error
	Update(ctx context.Context, w io.Writer, r io.Reader) error
	Delete(ctx context.Context, w io.Writer, r io.Reader) error
}

type Client struct {
	AuthOperations
	EventOperations
}

type Error struct {
	Message string `json:"message"`
}

func NewClient(ao AuthOperations, eo EventOperations) *Client {
	return &Client{
		ao,
		eo,
	}
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
	err := c.EventOperations.Update(r.Context(), w, r.Body)
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
	err := c.EventOperations.Delete(r.Context(), w, nil)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-delete")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}
