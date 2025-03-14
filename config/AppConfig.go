package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort string
	Dsn        string
	Secret     string
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

	return AppConfig{
		ServerPort: httpPort,
		Dsn:        dsn,
		Secret:     secret,
	}, nil
}
