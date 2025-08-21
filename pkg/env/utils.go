package env

import (
	"os"
	"strconv"
)

func GetEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvIntOrDefault(key string, fallback int) int {
	if value := GetEnvOrDefault(key, ""); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
