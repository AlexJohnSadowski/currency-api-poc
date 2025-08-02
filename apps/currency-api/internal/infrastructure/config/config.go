package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                string
	GinMode             string
	LogLevel            string
	OpenExchangeAPIKey  string
	OpenExchangeBaseURL string
	RedisURL            string
	Environment         string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		GinMode:             getEnv("GIN_MODE", "debug"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		OpenExchangeAPIKey:  getEnv("OPEN_EXCHANGE_API_KEY", ""),
		OpenExchangeBaseURL: getEnv("OPEN_EXCHANGE_BASE_URL", "https://openexchangerates.org/api"),
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		Environment:         getEnv("ENV", "development"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	if c.GinMode != "debug" && c.GinMode != "release" && c.GinMode != "test" {
		return fmt.Errorf("GIN_MODE must be one of: debug, release, test")
	}

	if c.LogLevel == "" {
		return fmt.Errorf("LOG_LEVEL cannot be empty")
	}

	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT must be a valid number: %w", err)
	}

	return nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.GinMode == "release"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
