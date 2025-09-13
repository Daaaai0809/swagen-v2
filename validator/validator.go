package validator

import (
	"errors"
	"regexp"

	"github.com/Daaaai0809/swagen-v2/input"
)

type IInputValidator interface {
	Validator_Alphanumeric_Underscore() *input.ValidationFunc
	Validator_Alphanumeric_Underscore_Allow_Empty() *input.ValidationFunc
}

type InputValidator struct{}

func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

func (v *InputValidator) Validator_Alphanumeric_Underscore() *input.ValidationFunc {
	var validator input.ValidationFunc = func(input string) error {
		// NOTE: only utf-8 alphanumeric and underscore are allowed
		var validName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
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
		var validName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
		if input != "" && !validName.MatchString(input) {
			return errors.New("parameter name can only contain alphanumeric characters and underscores, and cannot start with a number")
		}
		return nil
	}

	return &validator
}
