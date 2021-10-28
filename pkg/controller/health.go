package controller

import (
	"net/http"

	"calendar.com/pkg/response"
)

type Client struct{}

func (Client) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	r := response.NewPrint()
	r.PrettyPrint(w, "Im alive")
}
