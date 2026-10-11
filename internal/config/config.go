package config

import (
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addr              string        `env:"ADDR" envDefault:":8080"`
	LogLevel          slog.Level    `env:"LOG_LEVEL" envDefault:"info"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	MaxHeaderBytes    int           `env:"MAX_HEADER_BYTES" envDefault:"1048576"`
	RateLimit         int           `env:"RATE_LIMIT" envDefault:"100"`
	RateLimitWindow   time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"1m"`
	// Short, because file names carry no version: a deploy changes a file but not its URL.
	StaticCacheControl string `env:"STATIC_CACHE_CONTROL" envDefault:"public, max-age=3600"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}
