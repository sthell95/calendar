package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"calendar.com/pkg/logger"
)

type Repository interface {
	Create(interface{}) error
	FindById(interface{}, uuid.UUID) error
	Find(interface{}) error
}

func (r Storage) Create(entity interface{}) error {
	return r.Gorm.Create(entity).Error
}

func (r Storage) FindById(entity interface{}, id uuid.UUID) error {
	return r.Gorm.First(entity, id).Error
}

func (r Storage) Find(entity interface{}) error {
	return r.Gorm.Take(entity).Error
}

func NewDB(ctx context.Context) *gorm.DB {
	db, err := gorm.Open(postgres.Open(viper.GetString("postgres.url")), &gorm.Config{})
	if err != nil {
		logger.NewLogger().Write(logger.Error, "Postgres url is invalid", "db")
		log.Fatalln(err)
	}

	go func() {
		<-ctx.Done()
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		fmt.Println("Close", sqlDB.Close())
	}()
	logger.NewLogger().Write(logger.Info, "Postgre connection created", "db")

	return db
}

func (r *Storage) Close() error {
	db, err := r.Gorm.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
