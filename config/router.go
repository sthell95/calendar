package config

import (
	"context"
	"net/http"

	"calendar.com/pkg/controller"

	"github.com/gorilla/mux"
)

type Handlers struct {
	*controller.Controller
}

func Run(ctx context.Context, r *mux.Router) error {
	server := &http.Server{Addr: ":8000", Handler: r}

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}

func (h Handlers) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health_checker", h.HealthHandler).Methods(http.MethodGet)

	r.HandleFunc("/users", h.SignIn).Methods(http.MethodPost)

	return r
}

func (h *Handlers) NewHandler(c controller.Controller) {
	h.Controller = &c
}
