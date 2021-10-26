package controller

import (
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Im alive")
}
