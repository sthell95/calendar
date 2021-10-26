package config

import (
	"calendar.com/pkg/controller"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve() {
	r := mux.NewRouter()

	r.HandleFunc("/health_checker", controller.HealthHandler).Methods(http.MethodGet)

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(fmt.Sprintf("[serve]: %v", err))
	}
}
