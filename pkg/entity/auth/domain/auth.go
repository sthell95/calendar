package domain

import "calendar.com/pkg/entity/user/domain"

type Authentication interface {
	Login()
}

type Token string

type JWTAuth struct {
	user domain.User
}
