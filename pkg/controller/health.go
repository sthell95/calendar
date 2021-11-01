package controller

import (
	"net/http"

	"calendar.com/pkg/entity/user"

	"calendar.com/pkg/response"
)

type Client struct {
	HealthRepository user.UserRepository
}

func (_ *Client) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	r := response.NewPrint()
	r.PrettyPrint(w, "Im alive")
}
