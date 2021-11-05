package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/service"
)

func TestController_SignIn(t *testing.T) {
	credentials := entity.Credentials{
		Login:    "login",
		Password: "testtest",
	}
	tests := []struct {
		name        string
		services    *service.Services
		mock        func(*service.MockAuthorization, entity.Credentials)
		wantMessage string
		wantCode    int
	}{
		{
			name: "Valid",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				mock.EXPECT().CheckCredentials(creds).Return(nil).AnyTimes()
				mock.EXPECT().GenerateToken(creds.Login).Return(&entity.AuthToken{
					Token:     "token",
					ExpiresAt: 1,
				}, nil)
			},
			wantMessage: `{"token":"token","expires_at":1}`,
			wantCode:    http.StatusOK,
		},
		{
			name: "Invalid credentials",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				e := errors.New("Invalid credentials")
				mock.EXPECT().CheckCredentials(creds).Return(e).AnyTimes()
			},
			wantMessage: `{"message":"Invalid credentials"}`,
			wantCode:    http.StatusBadRequest,
		},
		{
			name: "Failed generation token",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				mock.EXPECT().CheckCredentials(creds).Return(nil).AnyTimes()
				e := errors.New("Couldn't generate token")
				mock.EXPECT().GenerateToken(creds.Login).Return(nil, e)
			},
			wantMessage: `{"message":"Couldn't generate token"}`,
			wantCode:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := service.NewMockAuthorization(ctrl)
			tt.mock(mock, credentials)

			w := httptest.NewRecorder()
			requestBody, _ := json.Marshal(credentials)
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))

			c := &Controller{
				Services: &service.Services{Authorization: mock},
			}
			c.SignIn(w, r)
			responseBody, _ := io.ReadAll(w.Body)
			if w.Code != tt.wantCode {
				t.Errorf("SignIn() failed: Response code: got: %d\n expected %d", w.Code, tt.wantCode)
			}

			token := strings.Trim(string(responseBody), "\n")
			if token != tt.wantMessage {
				t.Errorf("SignIn() failed: Response code: got: %v\n expected %v", token, tt.wantMessage)
			}
		})
	}
}
