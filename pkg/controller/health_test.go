package controller

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	type args func(w http.ResponseWriter, r *http.Request)
	tests := []struct {
		name    string
		handler args
		want    string
	}{
		{
			name: "Check alive",
			handler: func(w http.ResponseWriter, r *http.Request) {
				HealthHandler(w, r)
			},
			want: `"Im alive"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health_checker", nil)
			w := httptest.NewRecorder()
			test.handler(w, req)

			body, _ := io.ReadAll(w.Body)
			expect := strings.Trim(string(body), "\n")
			if expect != test.want {
				t.Errorf("Failed " + test.name + "\n Expected: " + test.want + "\n Got: " + string(body))
			}
		})
	}
}
