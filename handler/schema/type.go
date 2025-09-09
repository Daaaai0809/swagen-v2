package schema

import (
	"errors"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/handler"
	"github.com/Daaaai0809/swagen-v2/utils"
	"gopkg.in/yaml.v2"
)

type SchemaName string

type Schema struct {
	*handler.Property
	Input utils.IInputMethods
}

func NewSchema(input utils.IInputMethods) Schema {
	return Schema{
		Input:    input,
		Property: handler.NewProperty(input, "", nil, constants.MODE_SCHEMA),
	}
}

func (s Schema) InputSchemaName(name *SchemaName) error {
	var validate utils.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("[ERROR] schema name cannot be empty")
		}
		if len(input) > 100 {
			return errors.New("[ERROR] schema name cannot exceed 100 characters")
		}

		return nil
	}

	err := s.Input.StringInput((*string)(name), "Schema Name", &validate)
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
