package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		GRPC
		HTTP
		Postgres
		Outbox
		Logger
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT"`
		Host        string `env:"GRPC_HOST"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT"`
	}

	HTTP struct {
		UsePrefork bool `env:"HTTP_USE_PREFORK"`
	}

	Postgres struct {
		PoolMax int    `env:"POSTGRES_POOL_MAX"`
		URL     string `env:"POSTGRES_URL"`
	}

	Outbox struct {
		Workers         int           `env:"OUTBOX_WORKERS"`
		BatchSize       int           `env:"OUTBOX_BATCH_SIZE"`
		WaitTimeMS      time.Duration `env:"OUTBOX_WAIT_TIME_MS"`
		InProgressTTLMS time.Duration `env:"OUTBOX_IN_PROGRESS_TTL_MS"`
	}

	Logger struct {
		LogFilePath string `env:"LOG_PATH"`
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
