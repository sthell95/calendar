package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"calendar.com/pkg/storage"

	"github.com/spf13/viper"

	"calendar.com/config"
	"calendar.com/pkg/controller"
	"calendar.com/pkg/domain/repository"
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
	db := storage.NewDB(ctx)

	handlers := new(config.Handlers)
	storages := storage.Storage{Gorm: db}
	repos := repository.NewRepository(&storages)
	services := service.NewService(repos)
	controller.NewController(services)

	err := config.Run(ctx, handlers.NewRouter())
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
