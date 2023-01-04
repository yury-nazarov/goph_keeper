package options

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Сервис можно запустить как с помощью ключей, так и используя переменные окружения

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DB                   string `env:"DATABASE_URI"`
	MigrationDownTo     string `env:"MIGRATION_DOWN_TO"`
}

func NewConfig() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "set server address, by example: 127.0.0.1:8080")
	flag.StringVar(&cfg.DB, "d", cfg.DB, "set database URI for Postgres, by example: host=localhost port=5432 user=example password=123 dbname=example sslmode=disable connect_timeout=5")
	flag.StringVar(&cfg.MigrationDownTo, "m", cfg.MigrationDownTo, "set up or down. By default Up. Example: down")

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	flag.Parse()
	return cfg, nil
}
