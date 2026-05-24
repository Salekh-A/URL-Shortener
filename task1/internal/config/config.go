package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr     string
	BaseURL  string
	LogLevel string
}

func New() *Config {
	var cfg Config
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "address to run server")
	flag.StringVar(&cfg.BaseURL, "b", "localhost:8080", "base URL for short links")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	}

	return &cfg
}
