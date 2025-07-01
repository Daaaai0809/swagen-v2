package utils

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
)

type ValidationFunc func(input string) error

type IInputMethods interface {
	StringInput(result *string, label string, validation *ValidationFunc) error
	IntInput(result *int, label string, validation *ValidationFunc) error
	Int64Input(result *int64, label string, validation *ValidationFunc) error
	UInt32Input(result *uint32, label string, validation *ValidationFunc) error
	UInt64Input(result *uint64, label string, validation *ValidationFunc) error
	Float32Input(result *float32, label string, validation *ValidationFunc) error
	Float64Input(result *float64, label string, validation *ValidationFunc) error
	BooleanInput(result *bool, label string) error
	SelectInput(result *string, label string, items []string) error
}

type InputMethods struct{}

func NewInputMethods() IInputMethods {
	return &InputMethods{}
}

func (im *InputMethods) StringInput(result *string, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	*result = input
	return nil
}

func (im *InputMethods) IntInput(result *int, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value int
	value, err = strconv.Atoi(input)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) UInt32Input(result *uint32, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value uint64
	value, err = strconv.ParseUint(input, 10, 32)
	if err != nil {
		return err
	}

	*result = uint32(value)
	return nil
}

func (im *InputMethods) Int64Input(result *int64, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value int64
	value, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) UInt64Input(result *uint64, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value uint64
	value, err = strconv.ParseUint(input, 10, 64)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) Float32Input(result *float32, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value float64
	value, err = strconv.ParseFloat(input, 32)
	if err != nil {
		return err
	}

	*result = float32(value)
	return nil
}

func (im *InputMethods) Float64Input(result *float64, label string, validation *ValidationFunc) error {
	var prompt promptui.Prompt

	if validation != nil {
		prompt = promptui.Prompt{
			Label: label,
		}
	} else {
		prompt = promptui.Prompt{
			Label:    label,
			Validate: promptui.ValidateFunc(*validation),
		}
	}

	input, err := prompt.Run()
	if err != nil {
		return err
	}

	var value float64
	value, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) BooleanInput(result *bool, label string) error {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"true", "false"},
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "{{ . | cyan }}",
			Inactive: "{{ . | faint }}",
			Selected: "{{ . | green }}",
		},
	}

	_, input, err := prompt.Run()
	if err != nil {
		return err
	}

	value, err := strconv.ParseBool(input)
	if err != nil {
		return err
	}

	*result = value
	return nil
}

func (im *InputMethods) SelectInput(result *string, label string, items []string) error {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "{{ . | cyan }}",
			Inactive: "{{ . | faint }}",
			Selected: "{{ . | green }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return err
	}

	if index < 0 || index >= len(items) {
		return errors.New("invalid selection index")
	}

	*result = items[index]
	return nil
}
