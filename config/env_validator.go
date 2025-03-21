package config

import (
	"os"
)

// Exits program if the environment variable isn't found
func getOrAppendMissing(key string, missing []string) (string, []string) {
	value := os.Getenv(key)

	if value == "" {
		missing = append(missing, key)
	}

	return value, missing
}
