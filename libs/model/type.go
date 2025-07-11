package model

import (
	"errors"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/libs"
	"github.com/Daaaai0809/swagen-v2/utils"
)

type Model struct {
	Input      utils.IInputMethods    `yaml:"-"`
	Title      string                 `yaml:"title,omitempty"`
	Type       string                 `yaml:"type"`
	Properties map[string]libs.Schema `yaml:"properties,omitempty"`
}

func NewModel(input utils.IInputMethods) *Model {
	return &Model{
		Input:      input,
		Title:      "",
		Type:       constants.OBJECT_TYPE,
		Properties: make(map[string]libs.Schema),
	}
}

func (m *Model) ReadTitle() error {
	var validate utils.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("[ERROR] title cannot be empty")
		}
		if len(input) > 100 {
			return errors.New("[ERROR] title cannot exceed 100 characters")
		}

		return nil
	}

	err := m.Input.StringInput(&m.Title, "Model Title", &validate)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) ReadPropertyName(name *string) error {
	var validate utils.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("[ERROR] property name cannot be empty")
		}
		if len(input) > 100 {
			return errors.New("[ERROR] property name cannot exceed 100 characters")
		}

		return nil
	}

	err := m.Input.StringInput(name, "Property Name", &validate)
	if err != nil {
		return err
	}

	return nil
}
