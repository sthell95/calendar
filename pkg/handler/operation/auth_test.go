package operation

import (
	"context"
	"errors"
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
		name           string
		authService    *service.Authorization
		mock           func(*service.MockAuthorization, entity.Credentials)
		errorMessage   string
		entityResponse *entity.AuthToken
	}{
		{
			name: "Valid",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				mock.EXPECT().SignInProcess(gomock.Any(), &creds).Return(&entity.AuthToken{
					Token:     "token",
					ExpiresAt: 1,
				}, nil).AnyTimes()
			},
			errorMessage: "",
			entityResponse: &entity.AuthToken{
				Token:     "token",
				ExpiresAt: 1,
			},
		},
		{
			name: "Invalid credentials",
			mock: func(mock *service.MockAuthorization, creds entity.Credentials) {
				e := errors.New("Invalid credentials")
				mock.EXPECT().SignInProcess(gomock.Any(), &creds).Return(nil, e).AnyTimes()
			},
			errorMessage:   "Invalid credentials",
			entityResponse: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := service.NewMockAuthorization(ctrl)
			tt.mock(mock, credentials)

			c := NewAuthOperations(mock)
			token, err := c.SignIn(context.Background(), &credentials)

			require.Equal(t, tt.errorMessage, err.Error())
			require.Equal(t, tt.entityResponse, token)
		})
	}
}
