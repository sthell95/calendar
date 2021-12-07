package handler

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	t.Run("Check health", func(t *testing.T) {
		w := httptest.NewRecorder()
		client := Controller{}
		client.HealthHandler(w, nil)
		want := `"Im alive"`

		body, _ := io.ReadAll(w.Body)
		expect := strings.Trim(string(body), "\n")

		require.Equal(t, expect, want)
	})
}
