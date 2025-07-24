package util

import "os"

// GetEnv retrieves an environment variable or returns a default value if not set
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}