package rest

import (
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
	"context"
	"io"
	"net/http"
)

type AuthOperations interface {
	signIn(ctx context.Context, w io.Writer, r io.Reader) error
}

type EventOperations interface {
	create(ctx context.Context, w io.Writer, r io.Reader) error
	update(ctx context.Context, w io.Writer, r io.Reader) error
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

func (c *Client) SignIn(w http.ResponseWriter, r *http.Request) {
	err := c.AuthOperations.signIn(r.Context(), w, r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}

func (c *Client) EventCreate(w http.ResponseWriter, r *http.Request) {
	err := c.EventOperations.create(r.Context(), w, r.Body)
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
	err := c.EventOperations.update(r.Context(), w, r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "event-update")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}
