package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Port            int           `env:"GRPC_PORT" env-default:"50051"`
	LogLevel        string        `env:"LOG_LEVEL" env-default:"info"`
	Timeout         time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`
	GwPort          int           `env:"GRPC_GATEWAY_PORT" env-default:"8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
