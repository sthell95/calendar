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

type AuthService struct {
	UserRepository repository.UserRepository
}

type Authorization interface {
	GenerateToken(login string) (*entity.AuthToken, error)
	CheckCredentials(entity.Credentials) error
}

func (AuthService) GenerateToken(login string) (*entity.AuthToken, error) {
	expiresAt := time.Now().Add(time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.CustomClaims{
		Login: login,
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
	if err != nil && u == nil {
		return err
	}
	return matchPasswords(c.Password, u.Password)
}

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashBytes), err
}

func matchPasswords(hashed, current string) error {
	return bcrypt.CompareHashAndPassword([]byte(current), []byte(hashed))
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{
		UserRepository: repository.NewUserRepository(repo.Storage),
	}
}
