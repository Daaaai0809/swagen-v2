package libs

import (
	"errors"
	"fmt"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/utils"
)

type Schema struct {
	Input        utils.IInputMethods `yaml:"-"`
	PropertyName string              `yaml:"-"`
	ParentSchema *Schema             `yaml:"-"` // Optional parent schema for nested properties
	Mode         constants.InputMode `yaml:"-"` // Mode of the schema (MODEL, SCHEMA, API)
	Type         string              `yaml:"type"`
	Format       string              `yaml:"format,omitempty"`
	Properties   map[string]Schema   `yaml:"properties,omitempty"`
	Required     []string            `yaml:"required,omitempty"`
	Nullable     bool                `yaml:"nullable,omitempty"`
	Items        *Schema             `yaml:"items,omitempty"`
}

func NewSchema(input utils.IInputMethods, propertyName string, parentSchema *Schema, mode constants.InputMode) Schema {
	return Schema{
		Input:        input,
		PropertyName: propertyName,
		ParentSchema: parentSchema,
		Mode:         mode,
		Type:         "",
		Format:       "",
		Properties:   make(map[string]Schema),
		Required:     []string{},
		Nullable:     false,
		Items:        nil,
	}
}

func (s *Schema) ReadType() error {
	err := s.Input.SelectInput(&s.Type, "Select Property Type", constants.FieldTypeList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Schema) ReadFormat() error {
	err := s.Input.SelectInput(&s.Format, "Select Property Format", constants.FormatList[s.Type])
	if err != nil {
		return err
	}

	return nil
}

func (s *Schema) ReadRequired() error {
	isRequired := false
	err := s.Input.BooleanInput(&isRequired, "Is this property required?")
	if err != nil {
		return err
	}

	if isRequired {
		s.ParentSchema.Required = append(s.ParentSchema.Required, s.PropertyName)
	}

	return nil
}

func (s *Schema) ReadNullable() error {
	err := s.Input.BooleanInput(&s.Nullable, "Is this property nullable?")
	if err != nil {
		return err
	}

	return nil
}

func (s *Schema) ReadProperty() error {
	var propertyName string

	var validate utils.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("[ERROR] property name cannot be empty")
		}
		if len(input) > 100 {
			return errors.New("[ERROR] property name cannot exceed 100 characters")
		}
		if _, exists := s.Properties[input]; exists {
			return errors.New("[ERROR] property name already exists")
		}
		return nil
	}

	if err := s.Input.StringInput(&propertyName, "Property Name", &validate); err != nil {
		return err
	}

	propertySchema := NewSchema(s.Input, propertyName, s, s.Mode)
	if err := propertySchema.ReadSchema(); err != nil {
		return err
	}

	s.Properties[propertyName] = propertySchema

	switch propertySchema.Type {
	case constants.ARRAY_TYPE:
		if err := s.ReadItem(); err != nil {
			return err
		}
	case constants.OBJECT_TYPE:
		if err := s.ReadProperty(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) ReadItem() error {
	if s.Type != constants.ARRAY_TYPE {
		return errors.New("[ERROR] items can only be defined for array type")
	}

	if s.Items == nil {
		s.Items = &Schema{
			Input:      s.Input,
			Type:       "",
			Format:     "",
			Properties: make(map[string]Schema),
			Required:   []string{},
			Nullable:   false,
			Items:      nil,
		}
	}

	if err := s.Items.ReadSchema(); err != nil {
		return err
	}

	return nil
}

func (s *Schema) ReadSchema() error {
	if err := s.ReadType(); err != nil {
		return err
	}

	isUseFormat := false
	if s.Type == constants.STRING_TYPE || s.Type == constants.NUMBER_TYPE || s.Type == constants.INTEGER_TYPE {
		err := s.Input.BooleanInput(&isUseFormat, "Do you want to specify a format for this property?")
		if err != nil {
			return err
		}
	}

	if isUseFormat {
		if err := s.ReadFormat(); err != nil {
			return err
		}
	}

	if s.Mode != constants.MODE_MODEL {
		if err := s.ReadRequired(); err != nil {
			return err
		}
	}

	if err := s.ReadNullable(); err != nil {
		return err
	}

	if s.Type == constants.OBJECT_TYPE {
		for {
			if err := s.ReadProperty(); err != nil {
				return err
			}

			msg := fmt.Sprintf("Do you want to add another property? (%s)", s.PropertyName)

			isAdd := false
			err := s.Input.BooleanInput(&isAdd, msg)
			if err != nil {
				return err
			}

			if !isAdd {
				break
			}
		}
	} else if s.Type == constants.ARRAY_TYPE {
		if err := s.ReadItem(); err != nil {
			return err
		}
	}

	return nil
}
