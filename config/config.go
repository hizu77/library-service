package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		GRPC
		HTTP
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT"`
		Host        string `env:"GRPC_HOST"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT"`
	}

	HTTP struct {
		UsePrefork bool `env:"HTTP_USE_PREFORK"`
	}
)

var (
	ErrInvalidGrpcPort    = errors.New("invalid grpc port")
	ErrInvalidGatewayPort = errors.New("invalid gateway port")
	ErrInvalidGatewayHost = errors.New("invalid gateway host")
)

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.GRPC.Port = os.Getenv("GRPC_PORT")

	if cfg.GRPC.Port == "" {
		return nil, ErrInvalidGrpcPort
	}

	cfg.GRPC.Host = os.Getenv("GRPC_HOST")
	if cfg.GRPC.Host == "" {
		return nil, ErrInvalidGatewayHost
	}

	cfg.GRPC.GatewayPort = os.Getenv("GRPC_GATEWAY_PORT")

	if cfg.GRPC.GatewayPort == "" {
		return nil, ErrInvalidGatewayPort
	}

	cfg.HTTP.UsePrefork = os.Getenv("HTTP_USE_PREFORK") == "true"

	return cfg, nil
}
