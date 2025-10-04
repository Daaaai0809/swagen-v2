package validator

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Daaaai0809/swagen-v2/input"
)

type IInputValidator interface {
	Validate_Environment_Props() error
	Validator_Alphanumeric_Underscore() *input.ValidationFunc
	Validator_Alphanumeric_Underscore_Allow_Empty() *input.ValidationFunc
}

type InputValidator struct{}

func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

func (v *InputValidator) Validate_Environment_Props() error {
	requiredVars := []string{"SWAGEN_MODEL_PATH", "SWAGEN_SCHEMA_PATH", "SWAGEN_API_PATH"}
	for _, envVar := range requiredVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("required environment variable %s is not set", envVar)
		}
	}
	return nil
}

func (v *InputValidator) Validator_Alphanumeric_Underscore() *input.ValidationFunc {
	var validator input.ValidationFunc = func(input string) error {
		// NOTE: only utf-8 alphanumeric and underscore are allowed
		validName := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
		if !validName.MatchString(input) {
			return errors.New("parameter name can only contain alphanumeric characters and underscores, and cannot start with a number")
		}
		return nil
	}

	return &validator
}

func (v *InputValidator) Validator_Alphanumeric_Underscore_Allow_Empty() *input.ValidationFunc {
	var validator input.ValidationFunc = func(input string) error {
		// NOTE: only utf-8 alphanumeric and underscore are allowed
		validName := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
		if input != "" && !validName.MatchString(input) {
			return errors.New("parameter name can only contain alphanumeric characters and underscores, and cannot start with a number")
		}
		return nil
	}

	return &validator
}
