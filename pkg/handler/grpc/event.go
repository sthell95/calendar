package grpc

import (
	"context"
	"io"

	pg "calendar.com/proto"
)

type EventOperations interface {
	Create(ctx context.Context, w io.Writer, r io.Reader) error
	Update(ctx context.Context, w io.Writer, r io.Reader, eventId string) error
	Delete(ctx context.Context, w io.Writer, eventId string) error
}

type AuthHandler struct {
	pg.UnimplementedAuthServiceServer
	AuthOperations AuthOperations
}

func (h *AuthHandler) CreateEvent(ctx context.Context) {

}
