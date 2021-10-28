package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"calendar.com/config"
	"calendar.com/pkg/logger"
)

func main() {
	fmt.Println("-> Running application")

	err := config.Serve()
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}

	logrus.Info("qwerty")
}
