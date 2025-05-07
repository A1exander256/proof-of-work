package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App    App
	Server Server
	Pow    Pow
}

type App struct {
	LogLevel string `envconfig:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`
}

type Server struct {
	Host      string        `envconfig:"SERVER_HOST"                   validate:"required"`
	Port      int           `envconfig:"SERVER_PORT"                   validate:"required"`
	KeepAlive time.Duration `envconfig:"SERVER_KEEP_ALIVE,default=10s" validate:"min=10s"`
	Deadline  time.Duration `envconfig:"SERVER_DEADLINE,default=10s"   validate:"min=1s"`
}

type Pow struct {
	Difficulty uint8 `envconfig:"POW_DIFFICULTY" validate:"min=10"`
}

func Parse() (Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("reading .env file: %w", err)
	}

	if err := envconfig.Process("", cfg); err != nil {
		return cfg, fmt.Errorf("reading environments: %w", err)
	}

	if err := validator.New().Struct(cfg); err != nil {
		return cfg, fmt.Errorf("validation error: %w", err)
	}

	return cfg, nil
}
