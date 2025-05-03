package service

import (
	"github.com/reminders/internal/redis"
)

func SetRedisKey(key string, value interface{}) {
	err := redis.Client.Set(redis.Ctx, key, value, 3600).Err()
	if err != nil {
		panic(err)
	}
}
