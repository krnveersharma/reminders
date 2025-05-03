package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

type redisClient struct {
	REDIS_ADDR string
	REDIS_PASS string
	Ctx        context.Context
}

func SetupRedis(redisAddress, redisPassword string) *redisClient {
	return &redisClient{
		REDIS_ADDR: redisAddress,
		REDIS_PASS: redisPassword,
		Ctx:        context.Background(),
	}
}

func (r *redisClient) InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     r.REDIS_ADDR,
		Password: r.REDIS_PASS,
		DB:       0,
	})

	if err := Client.Ping(r.Ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis %v", err)
	}
}
