package entity

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Login    string
	Password string
}

type AuthToken struct {
	Token     string
	ExpiresAt int64
}

type CustomClaims struct {
	UserId uuid.UUID
	jwt.StandardClaims
}
