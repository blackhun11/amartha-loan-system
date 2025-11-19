package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App App `envconfig:"APP"`
}

type App struct {
	ServerPort string `envconfig:"SERVER_PORT" default:"1323"`
}

var instance Config

func Load() {
	var envFilePath = ".env"
	if _, err := os.Stat(envFilePath); err == nil {
		if err := godotenv.Load(envFilePath); err != nil {
			panic(err)
		}
	}
	err := envconfig.Process("", &instance)
	if err != nil {
		panic(err)
	}

}

func Instance() Config {
	return instance
}
