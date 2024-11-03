package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type key string

const KeyUUID key = key("uuid")
const KeyMetrics = key("metrics")

type Config struct {
	Service  Service
	Metrics  Metrics
	Platform Platform
}

type Service struct {
	Port string `env:"SEARCH_SERVICE_PORT"`
	Host string `env:"SEARCH_SERVICE_HOST"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
