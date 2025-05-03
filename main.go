package main

import (
	"log"

	"github.com/reminders/config"
	"github.com/reminders/internal/api"
	"github.com/reminders/internal/redis"
)

func main() {
	config, err := config.SetupEnv()
	if err != nil {
		log.Fatalf(err.Error())
	}

	setupRedis := redis.SetupRedis(config.RedisAddress, config.RedisPassword)
	setupRedis.InitRedis()

	api.StartServer(config)
}
