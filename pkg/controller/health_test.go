package controller

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	t.Run("Check health", func(t *testing.T) {
		w := httptest.NewRecorder()
		client := Controller{}
		client.HealthHandler(w, nil)
		want := `"Im alive"`

		body, _ := io.ReadAll(w.Body)
		expect := strings.Trim(string(body), "\n")
		if expect != want {
			t.Errorf("Failed Check health\n Expected: " + want + "\n Got: " + string(body))
		}
	})
}
