package options

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Сервис можно запустить как с помощью ключей, так и используя переменные окружения

type Config struct {
	Host        string `env:"HOST"`
	Port        int    `env:"PORT"`
	DB          string `env:"DATABASE_URI"`
	MigrateTo   string `env:"MIGRATE_TO"`
	MigrateFile string `env:"MIGRATE_FILE"`
	Version 	string
}

func NewConfig() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.Host, "h", cfg.Host, "127.0.0.1")
	flag.IntVar(&cfg.Port, "p", cfg.Port, "8080")
	flag.StringVar(&cfg.DB, "d", cfg.DB, "host=localhost port=5432 user=example password=123 dbname=example sslmode=disable connect_timeout=5")
	flag.StringVar(&cfg.MigrateTo, "mt", cfg.MigrateTo, "02")
	flag.StringVar(&cfg.MigrateFile, "mf", cfg.MigrateFile, "./internal/migrations")

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	flag.Parse()
	return cfg, nil
}
