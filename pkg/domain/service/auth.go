package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
)

type InvalidCredentials struct{}

func (InvalidCredentials) Error() string {
	return "Credentials error: Invalid credentials"
}

type PasswordNotMatched struct{}

func (PasswordNotMatched) Error() string {
	return "Password doesn't match"
}

type Notfound struct{}

func (Notfound) Error() string {
	return "User not found"
}

type AuthService struct {
	UserRepository repository.UserRepository
}

type Credentials interface {
	GenerateToken(*entity.Credentials) (*entity.AuthToken, error)
	CheckCredentials(entity.Credentials) error
}

type Authorization interface {
	SignInProcess(c *entity.Credentials) (*entity.AuthToken, error)
}

func (s *AuthService) SignInProcess(c *entity.Credentials) (*entity.AuthToken, error) {
	creds := entity.Credentials{
		Login:    c.Login,
		Password: c.Password,
	}
	err := s.CheckCredentials(creds)
	if err != nil {
		return nil, InvalidCredentials{}
	}

	token, err := s.GenerateToken(&creds)
	if err != nil {
		return nil, PasswordNotMatched{}
	}
	return token, nil
}

func (AuthService) GenerateToken(credentials *entity.Credentials) (*entity.AuthToken, error) {
	expiresAt := time.Now().Add(time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.CustomClaims{
		Login: credentials.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
		},
	})

	ex, _ := os.Executable()
	privatKeyByte, err := os.ReadFile(ex + viper.GetString("security.private_key"))
	signedToken, err := token.SignedString(privatKeyByte)

	return &entity.AuthToken{
		Token:     signedToken,
		ExpiresAt: expiresAt,
	}, err
}

func (s AuthService) CheckCredentials(c entity.Credentials) error {
	u, err := s.UserRepository.FindOneBy(map[string]interface{}{
		"login": c.Login,
	})
	if err != nil || u == nil {
		return Notfound{}
	}
	return matchPasswords(c.Password, u.Password)
}

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashBytes), err
}

func matchPasswords(hashed, current string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(current), []byte(hashed)); err != nil {
		return PasswordNotMatched{}
	}
	return nil
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: repo,
	}
}
