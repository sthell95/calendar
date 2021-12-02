package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"calendar.com/pkg/storage/mongodb"

	"calendar.com/pkg/domain/repository"

	"calendar.com/pkg/storage/postgresdb"

	"github.com/spf13/viper"

	"calendar.com/config"
	"calendar.com/pkg/controller"
	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
)

func main() {
	fmt.Println("-> Running application")

	if err := initConfig(); err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "main")
	}

	grace := make(chan os.Signal, 1)
	signal.Notify(grace, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		logger.NewLogger().Write(logger.Error, fmt.Sprintf("system call: %+v", <-grace), "main")
		cancel()
	}()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongo_url")))
	if err != nil {
		logger.NewLogger().Write(logger.Error, "Mongo url is invalid", "db")
		log.Fatalln(err)
	}

	gormClient, err := gorm.Open(postgres.Open(viper.GetString("postgresql_url")), &gorm.Config{})
	if err != nil {
		logger.NewLogger().Write(logger.Error, "Postgres url is invalid", "db")
		log.Fatalln(err)
	}

	eventMongoClient := mongodb.NewEventRepository(mongoClient)
	eventPostgresClient := postgresdb.NewRepository(gormClient)
	storageEventClient := repository.NewEventRepository(eventMongoClient, eventPostgresClient)

	eventService := service.NewEventService(storageEventClient)

	authService := service.NewAuthService(eventPostgresClient)
	c := controller.NewController(eventService, authService)
	handlers := new(config.Handlers)
	handlers.NewHandler(*c)

	err = config.Run(ctx, handlers.NewRouter())
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
