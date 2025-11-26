package config

import "github.com/pyr33x/envy"

type Config struct {
	Server Server
	Redis  Redis
	Zap    Zap
}

type Server struct {
	Origin Origin
	Proxy  Proxy
}

type Origin struct {
	URL string
}

type Proxy struct {
	Port string
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
		Server: Server{
			Origin: Origin{
				URL: envy.GetString("ORIGIN_URL", "http://dummyjson.com"),
			},
			Proxy: Proxy{
				Port: envy.GetString("PROXY_PORT", "1337"),
			},
		},
		Zap: Zap{
			Environment: envy.GetString("ZAP_ENV", "dev"),
		},
		Redis: Redis{
			Base: Base{
				Host:     envy.GetString("REDIS_HOST", "proxio_redis"),
				Port:     envy.GetString("REDIS_PORT", "6379"),
				Username: envy.GetString("REDIS_USERNAME", "proxio"),
				Password: envy.GetString("REDIS_PASSWORD", ""),
			},
			Database: envy.GetInt("REDIS_DATABASE", 0),
		},
	}
}
