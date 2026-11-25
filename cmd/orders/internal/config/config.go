package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ProductConfig struct {
	BaseURL string
}

type HTTPConfig struct {
	Timeout    time.Duration
	MaxRetries int
}

type Config struct {
	Products ProductConfig
	HTTP     HTTPConfig
}

func Load() (Config, error) {
	baseURL := os.Getenv("PRODUCTS_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	timeout := 5 * time.Second
	if value := os.Getenv("PRODUCTS_TIMEOUT"); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil || parsed <= 0 {
			return Config{}, fmt.Errorf("PRODUCTS_TIMEOUT must be a positive duration, for example 5s")
		}
		timeout = parsed
	}
	retries := 2
	if value := os.Getenv("HTTP_MAX_RETRIES"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 || parsed > 5 {
			return Config{}, fmt.Errorf("HTTP_MAX_RETRIES must be between 0 and 5")
		}
		retries = parsed
	}
	return Config{Products: ProductConfig{BaseURL: baseURL}, HTTP: HTTPConfig{Timeout: timeout, MaxRetries: retries}}, nil
}
