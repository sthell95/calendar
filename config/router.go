package config

import (
	"context"
	"net/http"

	"calendar.com/pkg/controller"

	"github.com/gorilla/mux"

	"calendar.com/pkg/entity/user"
)

type Handlers interface {
	HealthHandler(http.ResponseWriter, *http.Request)
}

func Serve(ctx context.Context, repo user.UserRepository) error {
	r := NewRouter(&controller.Client{HealthRepository: repo})
	server := &http.Server{Addr: ":8000", Handler: r}

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}

func NewRouter(h Handlers) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health_checker", h.HealthHandler).Methods(http.MethodGet)

	return r
}
