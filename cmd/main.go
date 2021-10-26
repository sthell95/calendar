package main

import (
	"fmt"

	"calendar.com/config"
)

func main() {
	fmt.Println("-> Running application")

	config.Serve()
}
