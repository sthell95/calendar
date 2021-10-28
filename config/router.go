package config

import (
	"net/http"

	"calendar.com/pkg/controller"

	"github.com/gorilla/mux"
)

func Serve() error {
	r := NewRouter()

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		return err
	}
	return nil
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health_checker", controller.HealthHandler).Methods(http.MethodGet)

	return r
}
