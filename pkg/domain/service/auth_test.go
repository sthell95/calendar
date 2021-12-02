package service

import (
	"errors"
	"testing"

	"calendar.com/pkg/storage/postgresdb"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"

	"calendar.com/pkg/domain/entity"
)

func TestAuthService_CheckCredentials(t *testing.T) {
	tests := []struct {
		name        string
		args        entity.Credentials
		condition   map[string]interface{}
		gotFromDb   *entity.User
		wantMessage error
		wantStruct  *entity.User
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
			wantStruct: &entity.User{
				Login:    "",
				Password: "",
			},
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
			wantMessage: errors.New("Password doesn't match"),
		},
		{
			name: "Client not found",
			args: entity.Credentials{
				Login:    "test",
				Password: "$2a$04$ESxZs3B48bQwdoWs03A8w.uVgiBZaHAC5Hoj1me9Ru0V/zFM4XIDG",
			},
			condition: map[string]interface{}{
				"login": "test",
			},
			gotFromDb:   nil,
			wantMessage: errors.New("Client not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepository := postgresdb.NewMockUserRepository(ctrl)
			mockUserRepository.EXPECT().FindOneBy(tt.condition).Return(tt.gotFromDb, nil).AnyTimes()
			s := AuthService{
				UserRepository: mockUserRepository,
			}

			_, err := s.CheckCredentials(tt.args)
			if err != nil {
				require.EqualError(t, err, tt.wantMessage.Error())
			}
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		credentials *entity.User
		wantErr     error
		config      func()
	}{
		{
			name:        "Valid generate token by credentials",
			credentials: &entity.User{Login: "test"},
			wantErr:     nil,
			config: func() {
				viper.Set("security.private_key", "./data/test.jwt.key")
				viper.Set("security.public_key", "./data/test.jwt.key.pub")
			},
		},
		{
			name:        "JWt keys don't exists",
			credentials: &entity.User{Login: "test"},
			wantErr:     errors.New("Could not generate token"),
			config:      func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config()
			ctrl := gomock.NewController(t)
			mock := postgresdb.NewMockUserRepository(ctrl)
			au := AuthService{
				UserRepository: mock,
			}
			got, err := au.GenerateToken(tt.credentials)
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
