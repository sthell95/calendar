package rest

import (
	"context"
	"io"
	"net/http"

	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type AuthOperations interface {
	SignIn(ctx context.Context, w io.Writer, r io.Reader) error
}

func (c *Client) SignIn(w http.ResponseWriter, r *http.Request) {
	err := c.AuthOperations.SignIn(r.Context(), w, r.Body)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "sign-in")
		response.NewPrint().PrettyPrint(
			w,
			Error{Message: err.Error()},
			response.WithCode(http.StatusBadRequest),
		)
	}
}
