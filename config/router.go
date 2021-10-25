package config

import (
	"fmt"
	"net/http"

	"calendar.com/handler"
	"github.com/gorilla/mux"
)

func Serve() {
	r := mux.NewRouter()

	r.HandleFunc("/health_checker", handler.HealthHandler).Methods(http.MethodGet)

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(fmt.Sprintf("[serve]: %v", err))
	}
}
