package handler

import (
	"errors"
	"fmt"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
)

type Property struct {
	Input          input.IInputMethods `yaml:"-"`
	PropertyName   string              `yaml:"-"`
	ParentProperty *Property           `yaml:"-"` // Optional parent schema for nested properties
	Mode           constants.InputMode `yaml:"-"` // Mode of the schema (MODEL, SCHEMA, API)

	Type       string               `yaml:"type,omitempty"`
	Format     string               `yaml:"format,omitempty"`
	Properties map[string]*Property `yaml:"properties,omitempty"`
	Required   []string             `yaml:"required,omitempty"`
	Nullable   bool                 `yaml:"nullable,omitempty"`
	Items      *Property            `yaml:"items,omitempty"`
	Example    string               `yaml:"example,omitempty"`
	Ref        string               `yaml:"$ref,omitempty"` // Reference to another schema
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
	label := "Select Property Type (" + s.PropertyName + ")"
	err := s.Input.SelectInput(&s.Type, label, constants.FieldTypeList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Property) readFormat() error {
	var format string

	label := "Select Property Format (" + s.PropertyName + ")"
	err := s.Input.SelectInput(&format, label, constants.FormatList[s.Type])
	if err != nil {
		return err
	}

	if format == constants.FORMAT_NONE {
		return nil
	}

	s.Format = format

	return nil
}

func (s *Property) readRequired() error {
	isRequired := false

	label := "Is this property required? (" + s.PropertyName + ")"
	err := s.Input.BooleanInput(&isRequired, label)
	if err != nil {
		return err
	}

	if isRequired {
		s.ParentProperty.Required = append(s.ParentProperty.Required, s.PropertyName)
	}

	return nil
}

func (s *Property) readNullable() error {
	label := "Is this property nullable? (" + s.PropertyName + ")"
	err := s.Input.BooleanInput(&s.Nullable, label)
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

func (s *Property) readPropertyNames() error {
	var propNames []string
	if err := s.Input.MultipleStringInput(&propNames, "Enter property names", nil); err != nil {
		return err
	}

	for _, name := range propNames {
		s.Properties[name] = NewProperty(s.Input, name, s, s.Mode)
	}

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

	if s.Type == constants.STRING_TYPE || s.Type == constants.NUMBER_TYPE || s.Type == constants.INTEGER_TYPE {
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
		if err := s.readPropertyNames(); err != nil {
			return err
		}

		if len(s.Properties) == 0 {
			err := fmt.Sprintf("[ERROR] at least one property is required for object type (property: %s)", s.PropertyName)
			return errors.New(err)
		}
		for _, prop := range s.Properties {
			if err := prop.ReadAll(); err != nil {
				return err
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
