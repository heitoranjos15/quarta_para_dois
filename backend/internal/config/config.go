package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisTLS      bool   `envconfig:"REDIS_TLS" default:"false"`
	GitHubToken   string `envconfig:"GITHUB_TOKEN" default:""`
	Port          string `envconfig:"PORT" default:"8080"`
	Env           string `envconfig:"ENV" default:"development"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
