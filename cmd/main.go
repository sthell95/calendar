package main

import (
	"fmt"

	"calendar.com/pkg/logger"

	"calendar.com/config"
)

func main() {
	fmt.Println("-> Running application")

	err := config.Serve()
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}
}
