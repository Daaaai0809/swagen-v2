package validator

import (
	"fmt"
	"os"
)

type IPropsValidator interface {
	Validate_Environment_Props() error
}

type PropsValidator struct{}

func NewPropsValidator() IPropsValidator {
	return &PropsValidator{}
}

func (v *PropsValidator) Validate_Environment_Props() error {
	requiredVars := []string{"SWAGEN_MODEL_PATH", "SWAGEN_SCHEMA_PATH", "SWAGEN_API_PATH"}
	for _, envVar := range requiredVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("required environment variable %s is not set", envVar)
		}
	}
	return nil
}
