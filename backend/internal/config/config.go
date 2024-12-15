package config

import (
	"fmt"

	"github.com/LLIEPJIOK/weather-forecast/backend/pkg/database/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.PostgresConfig
	RESTServerPort int `env:"REST_SERVER_PORT" env-default:"8080"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadConfig("./.env", &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &cfg, nil
}
