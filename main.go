/*
Copyright © 2025 NAME HERE dai.tsuruga0809@gmail.com
*/
package main

import (
	"fmt"
	"os"

	"github.com/Daaaai0809/swagen-v2/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	if err := validateEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	cmd.Execute()
}

func validateEnv() error {
	requiredVars := []string{"SWAGEN_MODEL_PATH", "SWAGEN_SCHEMA_PATH", "SWAGEN_API_PATH"}
	for _, envVar := range requiredVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("required environment variable %s is not set", envVar)
		}
	}
	return nil
}
