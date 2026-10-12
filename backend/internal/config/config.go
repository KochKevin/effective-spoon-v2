package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	StripeClientKey string `env:"STRIPE_PUBLIC_CLIENT_KEY"`
}

func ParseConfig() (Config, error) {

	//Load Env
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	//Parse Env
	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
