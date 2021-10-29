package domain

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	login    string
	password string
	timezone string
}
