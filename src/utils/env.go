package utils

import (
	"fmt"
	"os"
)

// GetEnvOrPanic panics if env variable not set
func GetEnvOrPanic(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required env var: %s", key))
	}
	return val
}

// GetEnvOrDefault returns fallback if env variable is empty
func GetEnvOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func GetEnvAsDurationInSeconds(key string, fallback string) int {
	val := GetEnvOrDefault(key, fallback)
	duration := parseDuration(val)
	return int(duration.Seconds())
}

func IsProduction() bool {
	return GetEnvOrDefault("APP_ENV", "development") == "production"
}
