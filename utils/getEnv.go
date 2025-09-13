package utils

import (
	"os"
)

const (
	SWAGEN_MODEL_PATH  = "SWAGEN_MODEL_PATH"
	SWAGEN_SCHEMA_PATH = "SWAGEN_SCHEMA_PATH"
	SWAGEN_API_PATH    = "SWAGEN_API_PATH"
)

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
