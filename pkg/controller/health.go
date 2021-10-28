package controller

import (
	"net/http"

	"calendar.com/pkg/response"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	r := response.NewPrint()
	r.PrettyPrint(w, "Im alive")
}
