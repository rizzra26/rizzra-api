package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgrespw!"`
	DBName     string `env:"DB_NAME" envDefault:"rizzra_dev"`
	JWTSecret  string `env:"JWT_SECRET" envDefault:"change-me-in-production"`
	Port       int    `env:"PORT" envDefault:"8888"`
	UploadDir  string `env:"UPLOAD_DIR" envDefault:"./uploads"`
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
