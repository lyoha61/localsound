package config

import (
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port                 string        `env:"PORT" envDefault:"8080"`
	DatabaseURL          string        `env:"DATABASE_URL" envDefault:"postgres://localsound:password@localhost:5432/localsound?sslmode=disable"`
	CookiesFromBrowser   bool          `env:"YOUTUBE_COOKIES_FROM_BROWSER" envDefault:"true"`
	CookiesPath          string        `env:"YOUTUBE_COOKIES_PATH" envDefault:"/app/cookies/youtube.txt"`
	StorageDir           string        `env:"STORAGE_DIR" envDefault:"storage/songs"`
	LrclibBaseURL        string        `env:"LRCLIB_BASE_URL" envDefault:"https://lrclib.net/api"`
	LrclibWorkerInterval time.Duration `env:"LRCLIB_WORKER_INTERVAL" envDefault:"10s"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Поддержка строк JDBC вида jdbc:postgresql://... для обратной совместимости со Spring конфигурацией
	if strings.HasPrefix(cfg.DatabaseURL, "jdbc:") {
		cfg.DatabaseURL = strings.TrimPrefix(cfg.DatabaseURL, "jdbc:")
	}

	return cfg, nil
}
