package postgres

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"calendar.com/pkg/logger"
)

type Repository interface {
	Create(interface{}) error
}

func (r *DB) Create(entity interface{}) error {
	return r.db.Create(entity).Error
}

func NewDB(ctx context.Context, url string) *DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
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
	return &DB{db: db}
}

func (r *DB) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
