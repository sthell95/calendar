package storage

import "gorm.io/gorm"

type Storage struct {
	Gorm *gorm.DB
}
