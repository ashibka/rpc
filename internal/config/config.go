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
	RedisHost       string        `env:"REDIS_HOST" env-default:"localhost"`
	RedisPort       string        `env:"REDIS_PORT" env-default:"6379"`
	RedisPassword   string        `env:"REDIS_PASSWORD" env-default:""`
	DbName          string        `env:"POSTGRES_DB" env-default:"postgres"`
	DbUser          string        `env:"POSTGRES_USER" env-default:"postgres"`
	DbPass          string        `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DbHost          string        `env:"POSTGRES_HOST" env-default:"db"`
	DbPort          int           `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresVersion string        `env:"POSTGRES_VERSION" env-default:"15"`
}

func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
