package entity

import "github.com/golang-jwt/jwt"

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type CustomClaims struct {
	Login string
	jwt.StandardClaims
}
