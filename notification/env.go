package notification

import (
	"github.com/Tsisar/extended-log-go/log"
	"os"
)

func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value == "true"
}
