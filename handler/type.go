package handler

import (
	"errors"
	"fmt"
	"slices"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/input"
)

type Property struct {
	Input              input.IInputMethods  `yaml:"-"`
	PropertyName       string               `yaml:"-"`
	ParentProperty     *Property            `yaml:"-"` // Optional parent schema for nested properties
	Mode               constants.InputMode  `yaml:"-"` // Mode of the schema (MODEL, SCHEMA, API)
	OptionalProperties *Optionals           `yaml:"-"`
	FileFetcher        fetcher.IFileFetcher `yaml:"-"`
	DirectoryPath      string               `yaml:"-"`

	Type       string               `yaml:"type,omitempty"`
	Format     string               `yaml:"format,omitempty"`
	Properties map[string]*Property `yaml:"properties,omitempty"`
	Required   []string             `yaml:"required,omitempty"`
	Nullable   bool                 `yaml:"nullable,omitempty"`
	Items      *Property            `yaml:"items,omitempty"`
	Example    string               `yaml:"example,omitempty"`
	Ref        string               `yaml:"$ref,omitempty"` // Reference to another schema
}

func NewProperty(input input.IInputMethods, propertyName string, parentProperty *Property, optionalProperties *Optionals, mode constants.InputMode, fileFetcher fetcher.IFileFetcher, directoryPath string) *Property {
	return &Property{
		Input:              input,
		PropertyName:       propertyName,
		ParentProperty:     parentProperty,
		Mode:               mode,
		DirectoryPath:      directoryPath,
		FileFetcher:        fileFetcher,
		OptionalProperties: optionalProperties,
		Type:               "",
		Format:             "",
		Properties:         make(map[string]*Property),
		Required:           []string{},
		Nullable:           false,
		Items:              nil,
		Example:            "",
	}
}

func (p *Property) readType() error {
	label := "Select Property Type (" + p.PropertyName + ")"
	err := p.Input.SelectInput(&p.Type, label, constants.FieldTypeList)
	if err != nil {
		return err
	}

	return nil
}

func (p *Property) readFormat() error {
	var format string

	label := "Select Property Format (" + p.PropertyName + ")"
	err := p.Input.SelectInput(&format, label, constants.FormatList[p.Type])
	if err != nil {
		return err
	}

	if format == constants.FORMAT_NONE {
		return nil
	}

	p.Format = format

	return nil
}

func (p *Property) readRequired() error {
	isRequired := false

	label := "Is this property required? (" + p.PropertyName + ")"
	err := p.Input.BooleanInput(&isRequired, label)
	if err != nil {
		return err
	}

	if isRequired {
		p.ParentProperty.Required = append(p.ParentProperty.Required, p.PropertyName)
	}

	return nil
}

func (p *Property) readNullable() error {
	label := "Is this property nullable? (" + p.PropertyName + ")"
	err := p.Input.BooleanInput(&p.Nullable, label)
	if err != nil {
		return err
	}

	return nil
}

func (p *Property) readRef() error {
	ref, err := p.FileFetcher.InteractiveResolveRef(p.Input, p.Mode, p.DirectoryPath)
	if err != nil {
		return err
	}
	if ref == "" {
		return errors.New("empty reference returned")
	}
	p.Ref = ref
	// when $ref is set, other siblings like type should be omitted in output
	p.Type = ""
	p.Format = ""
	p.Properties = nil
	p.Items = nil
	p.Nullable = false
	p.Example = ""
	return nil
}

// ReadProperty() is public method because this is used in schema handler
func (p *Property) ReadProperty() error {
	var propertyName string

	var validate input.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("[ERROR] property name cannot be empty")
		}
		if len(input) > 100 {
			return errors.New("[ERROR] property name cannot exceed 100 characters")
		}
		if _, exists := p.Properties[input]; exists {
			return errors.New("[ERROR] property name already exists")
		}
		return nil
	}

	if err := p.Input.StringInput(&propertyName, "Property Name", &validate); err != nil {
		return err
	}

	property := NewProperty(p.Input, propertyName, p, p.OptionalProperties, p.Mode, p.FileFetcher, p.DirectoryPath)

	if err := property.ReadAll(); err != nil {
		return err
	}

	p.Properties[propertyName] = property

	switch property.Type {
	case constants.ARRAY_TYPE:
		if err := p.ReadItem(); err != nil {
			return err
		}
	case constants.OBJECT_TYPE:
		if err := p.ReadProperty(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Property) ReadItem() error {
	if p.Type != constants.ARRAY_TYPE {
		return errors.New("[ERROR] items can only be defined for array type")
	}

	if p.Items == nil {
		p.Items = &Property{
			Input:      p.Input,
			Type:       "",
			Format:     "",
			Properties: make(map[string]*Property),
			Required:   []string{},
			Nullable:   false,
			Items:      nil,
			Example:    "",
		}
	}

	if err := p.Items.ReadProperty(); err != nil {
		return err
	}

	return nil
}

func (p *Property) readExample() error {
	var example string
	err := p.Input.StringInput(&example, "Example Value", nil)
	if err != nil {
		return err
	}
	p.Example = example
	return nil
}

func (p *Property) readPropertyNames() error {
	var propNames []string
	if err := p.Input.MultipleStringInput(&propNames, "Enter property names", nil); err != nil {
		return err
	}

	for _, name := range propNames {
		p.Properties[name] = NewProperty(p.Input, name, p, p.OptionalProperties, p.Mode, p.FileFetcher, p.DirectoryPath)
	}

	return nil
}

func (p *Property) ReadAll() error {
	if p.isReadRef() {
		var useRef bool

		label := "Do you want to reference another schema? (" + p.PropertyName + ")"
		if err := p.Input.BooleanInput(&useRef, label); err != nil {
			return err
		}

		if useRef {
			if err := p.readRef(); err != nil {
				return err
			}

			if err := p.readRequired(); err != nil {
				return err
			}

			return nil
		}
	}

	if err := p.readType(); err != nil {
		return err
	}

	if p.Type == constants.STRING_TYPE || p.Type == constants.NUMBER_TYPE || p.Type == constants.INTEGER_TYPE {
		if err := p.readFormat(); err != nil {
			return err
		}
	}

	if p.isReadRequired() {
		if err := p.readRequired(); err != nil {
			return err
		}
	}

	if p.isReadNullable() {
		if err := p.readNullable(); err != nil {
			return err
		}
	}

	switch p.Type {
	case constants.OBJECT_TYPE:
		if err := p.readPropertyNames(); err != nil {
			return err
		}

		if len(p.Properties) == 0 {
			err := fmt.Sprintf("[ERROR] at least one property is required for object type (property: %s)", p.PropertyName)
			return errors.New(err)
		}
		for _, prop := range p.Properties {
			if err := prop.ReadAll(); err != nil {
				return err
			}
		}
	case constants.ARRAY_TYPE:
		if err := p.ReadItem(); err != nil {
			return err
		}
	}

	if constants.IsExamplableType(p.Type) && p.OptionalProperties.Contains("example") {
		isAddExample := false
		err := p.Input.BooleanInput(&isAddExample, "Do you want to add an example value for this property?")
		if err != nil {
			return err
		}

		if isAddExample {
			if err := p.readExample(); err != nil {
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
	if p.Mode == constants.MODE_API {
		return true
	}

	return p.Mode != constants.MODE_MODEL && p.ParentProperty != nil
}

type Optionals []string

func (o Optionals) Contains(prop string) bool {
	return slices.Contains(o, prop)
}
