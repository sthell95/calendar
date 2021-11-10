package entity

import "github.com/golang-jwt/jwt"

type Credentials struct {
	Login    string
	Password string
}

type AuthToken struct {
	Token     string
	ExpiresAt int64
}

type CustomClaims struct {
	Login string
	jwt.StandardClaims
}

type key string

const CtxUserKey key = "user"
const CtxUserKey2 key = "user2"
