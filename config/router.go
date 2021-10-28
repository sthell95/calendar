package config

import (
	"net/http"

	"calendar.com/pkg/controller"

	"github.com/gorilla/mux"
)

type Handlers interface {
	HealthHandler(http.ResponseWriter, *http.Request)
}

func Serve() error {
	r := NewRouter(controller.Client{})

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		return err
	}
	return nil
}

func NewRouter(h Handlers) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health_checker", h.HealthHandler).Methods(http.MethodGet)

	return r
}
