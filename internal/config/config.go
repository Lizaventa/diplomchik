package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	PgxDBAddr  string `env:"PGXCONN" envDefault:"postgres://postgres:Kev0507_24@localhost:15432/lab"`
	SigningKey string `env:"SIGNING_KEY" envDefault:"fhrewbiyf234gbr2bgf742fg7635467fb2"`
}

// NewConfig creates a new Config instance
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
