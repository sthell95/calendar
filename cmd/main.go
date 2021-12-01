package main

import (
	"calendar.com/pkg/domain/repository/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/viper"

	"calendar.com/config"
	"calendar.com/pkg/controller"
	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/storage"
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

	storageClient := storage.NewClient(ctx)
	eventRepository := postgres.NewEventRepository(storageClient)
	eventService := service.NewEventService(eventRepository)
	userRepository := postgres.NewUserRepository(storageClient)
	authService := service.NewAuthService(userRepository)
	c := controller.NewController(eventService, authService)
	handlers := new(config.Handlers)
	handlers.NewHandler(*c)

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
