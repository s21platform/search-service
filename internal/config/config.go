package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Metrics  Metrics
	Platform Platform
	User     User
	Logger   Logger
	Society  Society
}

type Service struct {
	Port string `env:"SEARCH_SERVICE_PORT"`
	Name string `env:"SEARCH_SERVICE_NAME"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

type User struct {
	Host string `env:"USER_SERVICE_HOST"`
	Port string `env:"USER_SERVICE_PORT"`
}

type Logger struct {
	Host string `env:"LOGGER_SERVICE_HOST"`
	Port string `env:"LOGGER_SERVICE_PORT"`
}

type Society struct {
	Host string `env:"SOCIETY_SERVICE_HOST"`
	Port string `env:"SOCIETY_SERVICE_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
