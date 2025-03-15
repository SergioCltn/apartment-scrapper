package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	DBPath         string
	ScraperBaseURL string
	ScraperTimeout time.Duration
}

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		configInstance = &Config{
			DBPath:         getEnv("DB_PATH", "apartments.db"),
			ScraperBaseURL: getEnv("SCRAPER_BASE_URL", ""),
			ScraperTimeout: getEnvAsDuration("SCRAPER_TIMEOUT", 10*time.Second),
		}
	})
	return configInstance
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if ms, err := strconv.Atoi(value); err == nil {
			return time.Duration(ms) * time.Millisecond
		}
	}
	return defaultValue
}
