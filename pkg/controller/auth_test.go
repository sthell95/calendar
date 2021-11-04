package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calendar.com/pkg/domain/entity"

	"github.com/golang/mock/gomock"

	"calendar.com/pkg/domain/service"
)

func TestController_SignIn(t *testing.T) {
	tests := []struct {
		name     string
		services *service.Services
		args     *entity.Credentials
		mock     func(*service.MockAuthorization)
		want     string
	}{
		{
			name: "Valid",
			args: &entity.Credentials{
				Login:    "login",
				Password: "testtest",
			},
			mock: func(mock *service.MockAuthorization) {
				mock.EXPECT().CheckCredentials(entity.Credentials{}).Return()
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := service.NewMockAuthorization(ctrl)
			tt.mock(mock)

			w := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tt.args)
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))

			c := &Controller{
				Services: &service.Services{Authorization: mock},
			}
			c.SignIn(w, r)
			responseBody, _ := io.ReadAll(w.Body)

			if w.Code != tt.want.Code {
				t.Errorf("SignIn() failed: Response code: got: %w\n expected %w", w.Code, tt.want.Code)
			}

			token := strings.Trim(string(responseBody), "\n")

			if token != tt.want {
				t.Errorf("SignIn() failed: Response code: got: %w\n expected %w", w.Code, tt.want.Code)
			}
		})
	}
}
