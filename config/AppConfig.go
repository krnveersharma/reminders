package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort        string
	Dsn               string
	Secret            string
	RedisAddress      string
	RedisPassword     string
	RazorpayKeyId     string
	RazorpayKeySecret string
}

func SetupEnv() (cfg AppConfig, err error) {
	godotenv.Load()
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, errors.New("http port is not found")
	}

	dsn := os.Getenv("DSN")

	secret := os.Getenv("SECRET_KEY")
	if dsn == "" {
		return AppConfig{}, errors.New("DB details not found")
	}

	if secret == "" {
		return AppConfig{}, errors.New("no jwt token exist")
	}

	redisAddress := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	razorpayKeyId := os.Getenv("RAZORPAY_KEY_ID")
	razorpayKeySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	return AppConfig{
		ServerPort:        httpPort,
		Dsn:               dsn,
		Secret:            secret,
		RedisAddress:      redisAddress,
		RedisPassword:     redisPassword,
		RazorpayKeyId:     razorpayKeyId,
		RazorpayKeySecret: razorpayKeySecret,
	}, nil
}
