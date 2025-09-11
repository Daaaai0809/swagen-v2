package handler

import (
	"errors"
	"fmt"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
)

type Property struct {
	Input          input.IInputMethods  `yaml:"-"`
	PropertyName   string               `yaml:"-"`
	ParentProperty *Property            `yaml:"-"` // Optional parent schema for nested properties
	Mode           constants.InputMode  `yaml:"-"` // Mode of the schema (MODEL, SCHEMA, API)
	Type           string               `yaml:"type,omitempty"`
	Format         string               `yaml:"format,omitempty"`
	Properties     map[string]*Property `yaml:"properties,omitempty"`
	Required       []string             `yaml:"required,omitempty"`
	Nullable       bool                 `yaml:"nullable,omitempty"`
	Items          *Property            `yaml:"items,omitempty"`
	Example        string               `yaml:"example,omitempty"`
	Ref            string               `yaml:"$ref,omitempty"` // Reference to another schema
}

func NewProperty(input input.IInputMethods, propertyName string, parentProperty *Property, mode constants.InputMode) *Property {
	return &Property{
		Input:          input,
		PropertyName:   propertyName,
		ParentProperty: parentProperty,
		Mode:           mode,
		Type:           "",
		Format:         "",
		Properties:     make(map[string]*Property),
		Required:       []string{},
		Nullable:       false,
		Items:          nil,
		Example:        "",
	}
}

func (s *Property) readType() error {
	err := s.Input.SelectInput(&s.Type, "Select Property Type", constants.FieldTypeList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Property) readFormat() error {
	err := s.Input.SelectInput(&s.Format, "Select Property Format", constants.FormatList[s.Type])
	if err != nil {
		return err
	}

	return nil
}

func (s *Property) readRequired() error {
	isRequired := false
	err := s.Input.BooleanInput(&isRequired, "Is this property required?")
	if err != nil {
		return err
	}

	if isRequired {
		s.ParentProperty.Required = append(s.ParentProperty.Required, s.PropertyName)
	}

	return nil
}

func (s *Property) readNullable() error {
	err := s.Input.BooleanInput(&s.Nullable, "Is this property nullable?")
	if err != nil {
		return err
	}

	return nil
}

func (s *Property) readRef() error {
	ref, err := utils.InteractiveResolveRef(s.Input, s.Mode)
	if err != nil {
		return err
	}
	if ref == "" {
		return errors.New("empty reference returned")
	}
	s.Ref = ref
	// when $ref is set, other siblings like type should be omitted in output
	s.Type = ""
	s.Format = ""
	s.Properties = nil
	s.Items = nil
	s.Nullable = false
	s.Example = ""
	return nil
}

// ReadProperty() is public method because this is used in schema handler
func (s *Property) ReadProperty() error {
	var propertyName string

	var validate input.ValidationFunc = func(input string) error {
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

	property := NewProperty(s.Input, propertyName, s, s.Mode)

	if err := property.ReadAll(); err != nil {
		return err
	}

	s.Properties[propertyName] = property

	switch property.Type {
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

func (s *Property) ReadItem() error {
	if s.Type != constants.ARRAY_TYPE {
		return errors.New("[ERROR] items can only be defined for array type")
	}

	if s.Items == nil {
		s.Items = &Property{
			Input:      s.Input,
			Type:       "",
			Format:     "",
			Properties: make(map[string]*Property),
			Required:   []string{},
			Nullable:   false,
			Items:      nil,
			Example:    "",
		}
	}

	if err := s.Items.ReadProperty(); err != nil {
		return err
	}

	return nil
}

func (s *Property) readExample() error {
	var example string
	err := s.Input.StringInput(&example, "Example Value", nil)
	if err != nil {
		return err
	}
	s.Example = example
	return nil
}

func (s *Property) ReadAll() error {
	if s.isReadRef() {
		var useRef bool
		if err := s.Input.BooleanInput(&useRef, "Do you want to reference another schema?"); err != nil {
			return err
		}

		if useRef {
			if err := s.readRef(); err != nil {
				return err
			}
			return nil
		}
	}

	if err := s.readType(); err != nil {
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
		if err := s.readFormat(); err != nil {
			return err
		}
	}

	if s.isReadRequired() {
		if err := s.readRequired(); err != nil {
			return err
		}
	}

	if s.isReadNullable() {
		if err := s.readNullable(); err != nil {
			return err
		}
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

	if constants.IsExamplableType(s.Type) {
		isAddExample := false
		err := s.Input.BooleanInput(&isAddExample, "Do you want to add an example value for this property?")
		if err != nil {
			return err
		}

		if isAddExample {
			if err := s.readExample(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Property) isReadRequired() bool {
	if p.ParentProperty == nil {
		return false
	}

	if p.Mode == constants.MODE_MODEL {
		return false
	}

	return true
}

func (p *Property) isReadNullable() bool {
	return p.Mode != constants.MODE_API
}

func (p *Property) isReadRef() bool {
	if p.ParentProperty == nil {
		return true
	}

	return p.Mode != constants.MODE_MODEL
}
