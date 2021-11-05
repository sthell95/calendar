package service

import (
	"errors"
	"testing"

	"github.com/spf13/viper"

	"github.com/golang/mock/gomock"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
)

func TestAuthService_CheckCredentials(t *testing.T) {
	tests := []struct {
		name        string
		args        entity.Credentials
		condition   map[string]interface{}
		gotFromDb   *entity.User
		wantMessage error
	}{
		{
			name: "Valid",
			args: entity.Credentials{
				Login:    "test",
				Password: "testtest",
			},
			condition: map[string]interface{}{
				"login": "test",
			},
			gotFromDb: &entity.User{
				Login:    "test",
				Password: "$2a$04$ESxZs3B48bQwdoWs03A8w.uVgiBZaHAC5Hoj1me9Ru0V/zFM4XIDG",
			},
			wantMessage: nil,
		},
		{
			name: "Invalid password",
			args: entity.Credentials{
				Login:    "test",
				Password: "testtest1",
			},
			condition: map[string]interface{}{
				"login": "test",
			},
			gotFromDb: &entity.User{
				Login:    "test",
				Password: "$2a$04$ESxZs3B48bQwdoWs03A8w.uVgiBZaHAC5Hoj1me9Ru0V/zFM4XIDG",
			},
			wantMessage: errors.New("Authorization error: Passwords don't not matched"),
		},
		{
			name: "User not found",
			args: entity.Credentials{
				Login:    "test",
				Password: "$2a$04$ESxZs3B48bQwdoWs03A8w.uVgiBZaHAC5Hoj1me9Ru0V/zFM4XIDG",
			},
			condition: map[string]interface{}{
				"login": "test",
			},
			gotFromDb:   nil,
			wantMessage: errors.New("Authorization error: Invalid credentials"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := repository.NewMockUserRepository(ctrl)
			mockUserRepository.EXPECT().FindOneBy(tt.condition).Return(tt.gotFromDb, nil).AnyTimes()
			s := AuthService{
				UserRepository: mockUserRepository,
			}

			err := s.CheckCredentials(tt.args)
			if err != nil && tt.wantMessage != nil && errors.Is(err, tt.wantMessage) {
				t.Errorf("CheckCredentials() error = %v, wantMessage %v", err, tt.wantMessage)
			}
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr error
		config  func()
	}{
		{
			name:    "Valid generate token by login",
			login:   "test",
			wantErr: nil,
			config: func() {
				viper.Set("security.private_key", "./data/test.jwt.key")
				viper.Set("security.public_key", "./data/test.jwt.key.pub")
			},
		},
		{
			name:    "JWt keys don't exists",
			login:   "test",
			wantErr: errors.New("Could not generate token"),
			config:  func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config()
			ctrl := gomock.NewController(t)
			mock := repository.NewMockUserRepository(ctrl)
			au := AuthService{
				UserRepository: mock,
			}
			got, err := au.GenerateToken(tt.login)
			if err != tt.wantErr && got == nil {
				t.Errorf("GenerateToken() error = %v, wantMessage %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_hashPassword(t *testing.T) {
	t.Run("Hash password", func(t *testing.T) {
		_, err := hashPassword("some password")

		if err != nil {
			t.Errorf(err.Error())
		}
	})
}

func Test_matchPasswords(t *testing.T) {
	type args struct {
		requested string
		current   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Matched",
			args: args{
				requested: "testtest",
				current:   "$2a$04$ESxZs3B48bQwdoWs03A8w.uVgiBZaHAC5Hoj1me9Ru0V/zFM4XIDG",
			},
			wantErr: false,
		},
		{
			name: "Dont Matched",
			args: args{
				requested: "testtest123",
				current:   "$2a$04$ESxZs3B48bQwdoWs03A8w.uVgiBZaHAC5Hoj1me9Ru0V/zFM4XIDG",
			},
			wantErr: true,
		},
		{
			name: "Password in db is not hash",
			args: args{
				requested: "testtest",
				current:   "dfetewtew",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchPasswords(tt.args.requested, tt.args.current); (err != nil) != tt.wantErr {
				t.Errorf("matchPasswords() error = %v, wantMessage %v", err, tt.wantErr)
			}
		})
	}
}
