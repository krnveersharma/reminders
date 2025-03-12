package main

import (
	"log"

	"github.com/reminders/config"
	"github.com/reminders/internal/api"
)

func main() {
	config, err := config.SetupEnv()
	if err != nil {
		log.Fatalf(err.Error())
	}

	api.StartServer(config)
}
