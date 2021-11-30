package storage

import (
	"context"

	"github.com/spf13/viper"
)

func NewClient(ctx context.Context) Repository {
	client := viper.GetString("db_client")
	switch client {
	case "gorm":
		return NewDB(ctx)
	case "mongo":
		return NewDb(ctx)
	default:
		panic("You have to set valid db client in config")
	}
}
