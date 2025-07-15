package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		GRPC
		HTTP
		Postgres
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT, required"`
		Host        string `env:"GRPC_HOST, required"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT, required"`
	}

	HTTP struct {
		UsePrefork bool `env:"HTTP_USE_PREFORK, required"`
	}

	Postgres struct {
		PoolMax int    `env:"POSTGRES_POOL_MAX, required"`
		URL     string `env:"POSTGRES_URL, required"`
	}
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
