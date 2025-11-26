package config

import "github.com/pyr33x/envy"

type Config struct {
	Redis Redis
	Zap   Zap
}

type Zap struct {
	Environment string
}

type Redis struct {
	Base
	Database int
}

type Base struct {
	Host     string
	Port     string
	Username string
	Password string
}

func New() *Config {
	return &Config{
		Redis: Redis{
			Base: Base{
				Host:     envy.GetString("REDIS_HOST", "127.0.0.1"),
				Port:     envy.GetString("REDIS_PORT", "6379"),
				Username: envy.GetString("REDIS_USERNAME", "proxy"),
				Password: envy.GetString("REDIS_PASSWORD", ""),
			},
			Database: envy.GetInt("REDIS_DATABASE", 0),
		},
		Zap: Zap{
			Environment: envy.GetString("ZAP_ENV", "dev"),
		},
	}
}
