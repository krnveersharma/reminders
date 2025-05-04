package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/reminders/internal/redis"
)

func SetRedisKey(key string, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	err = redis.Client.Set(redis.Ctx, key, data, 3600*time.Second).Err()

	if err != nil {
		panic(err)
	}
	fmt.Printf("redis value set successfuly:%v\n", key)
}

func GetFromRedis(key string) (string, error) {
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
