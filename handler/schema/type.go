package schema

import (
	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/handler"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
	"gopkg.in/yaml.v2"
)

type SchemaName string

type Schema struct {
	*handler.Property
	Input     input.IInputMethods
	Validator validator.IInputValidator
}

func NewSchema(input input.IInputMethods, validator validator.IInputValidator) Schema {
	return Schema{
		Input:     input,
		Validator: validator,
		Property:  handler.NewProperty(input, "", nil, &handler.Optionals{}, constants.MODE_SCHEMA),
	}
}

func (s Schema) InputPropertyNames() error {
	for {
		var propName string
		if err := s.Input.StringInput(&propName, "Enter a property name", s.Validator.Validator_Alphanumeric_Underscore()); err != nil {
			return err
		}
		if propName == "" {
			break
		}

		s.Properties[propName] = handler.NewProperty(s.Input, propName, s.Property, s.OptionalProperties, constants.MODE_SCHEMA)
	}

	return nil
}

func (s Schema) InputSchemaName(name *SchemaName) error {
	err := s.Input.StringInput((*string)(name), "Schema Name", s.Validator.Validator_Alphanumeric_Underscore())
	if err != nil {
		return err
	}

	return nil
}

func (s *Schema) GenerateSchema(fileName string, schemaName SchemaName, path string) error {
	data, err := yaml.Marshal(map[SchemaName]*handler.Property{
		schemaName: s.Property,
	})
	if err != nil {
		return err
	}

	if err := utils.GenerateSchema(data, fileName, path); err != nil {
		return err
	}

	return nil
}
