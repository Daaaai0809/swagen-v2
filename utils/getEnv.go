package utils

import (
	"os"
)

const (
	MODEL_PATH  = "MODEL_PATH"
	SCHEMA_PATH = "SCHEMA_PATH"
	API_PATH    = "API_PATH"
)

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
