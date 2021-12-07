package handler

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
	"github.com/stretchr/testify/require"

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
		authService *service.Authorization
		mock        func(*service.MockAuthorization, entity.Credentials)
		wantMessage string
		wantCode    int
	}{
		{
			name: "Valid",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				mock.EXPECT().SignInProcess(gomock.Any(), &creds).Return(&entity.AuthToken{
					Token:     "token",
					ExpiresAt: 1,
				}, nil).AnyTimes()
			},
			wantMessage: `{"token":"token","expires_at":1}`,
			wantCode:    http.StatusOK,
		},
		{
			name: "Invalid credentials",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				e := errors.New("Invalid credentials")
				mock.EXPECT().SignInProcess(gomock.Any(), &creds).Return(nil, e).AnyTimes()
			},
			wantMessage: `{"message":"Invalid credentials"}`,
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
				EventService: nil,
				AuthService:  mock,
			}
			c.SignIn(w, r)
			responseBody, _ := io.ReadAll(w.Body)
			token := strings.Trim(string(responseBody), "\n")
			require.Equal(t, w.Code, tt.wantCode)
			require.Equal(t, token, tt.wantMessage)
		})
	}
}
