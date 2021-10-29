package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"calendar.com/config"
	"calendar.com/pkg/logger"
)

func main() {
	fmt.Println("-> Running repo")

	err := config.Serve()
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}

	logrus.Info("qwerty")
}
