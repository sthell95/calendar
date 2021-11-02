package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"calendar.com/pkg/entity/user"

	"calendar.com/pkg/storage/postgres"

	"calendar.com/config"
	"calendar.com/pkg/logger"
)

func main() {
	fmt.Println("-> Running application")

	grace := make(chan os.Signal, 1)
	signal.Notify(grace, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		logger.NewLogger().Write(logger.Error, fmt.Sprintf("system call: %+v", <-grace), "main")
		cancel()
	}()
	db := postgres.NewDB(ctx, os.Getenv("POSTGRE_URL"))
	userRepo := user.NewUserRepo(db)

	err := config.Serve(ctx, userRepo)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}
}
