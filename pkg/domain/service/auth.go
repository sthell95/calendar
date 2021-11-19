package service

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/domain/repository"
	"calendar.com/pkg/logger"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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

type NotAuthorized struct{}

func (NotAuthorized) Error() string {
	return "Not Authorized"
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
	h, _ := hashPassword(creds.Password)
	fmt.Println(h)
	user, err := s.CheckCredentials(creds)
	if err != nil {
		return nil, InvalidCredentials{}
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, PasswordNotMatched{}
	}
	return token, nil
}

func (AuthService) GenerateToken(u *entity.User) (*entity.AuthToken, error) {
	privatKeyByte, err := os.ReadFile(viper.GetString("security.private_key"))
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "token")
		return nil, err
	}

	expiresAt := time.Now().Add(time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.CustomClaims{
		UserId: u.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
		},
	})

	signedToken, err := token.SignedString(privatKeyByte)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "token")
		return nil, err
	}

	return &entity.AuthToken{
		Token:     signedToken,
		ExpiresAt: expiresAt,
	}, err
}

func (s AuthService) CheckCredentials(c entity.Credentials) (*entity.User, error) {
	u, err := s.UserRepository.FindOneBy(map[string]interface{}{
		"login": c.Login,
	})
	if err != nil || u == nil {
		return nil, Notfound{}
	}

	err = matchPasswords(c.Password, u.Password)
	if err != nil {
		return nil, err
	}

	return u, nil
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

func Validate(r *http.Request) (uuid.UUID, error) {
	gotToken := r.Header.Get("Authorization")
	t, err := jwt.ParseWithClaims(gotToken, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, NotAuthorized{}
		}

		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			return nil, fmt.Errorf("Expired token")
		}

		verifyBytes, err := os.ReadFile(viper.GetString("security.private_key"))
		if err != nil {
			return nil, err
		}

		return verifyBytes, nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := t.Claims.(*entity.CustomClaims)
	if !ok || !t.Valid || claims.UserId.String() == "" {
		return uuid.UUID{}, err
	}
	return claims.UserId, nil
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: repo,
	}
}
