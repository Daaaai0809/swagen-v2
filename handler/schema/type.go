package schema

import (
	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/handler"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
	"gopkg.in/yaml.v2"
)

type SchemaName string

type Schema struct {
	*handler.Property
	Input            input.IInputMethods
	Validator        validator.IInputValidator
	FileFetcher      fetcher.IFileFetcher
	DirectoryFetcher fetcher.IDirectoryFetcher
	DirectoryPath    string `yaml:"-"`
}

func NewSchema(input input.IInputMethods, validator validator.IInputValidator, fileFetcher fetcher.IFileFetcher, directoryFetcher fetcher.IDirectoryFetcher) Schema {
	return Schema{
		Input:            input,
		Validator:        validator,
		FileFetcher:      fileFetcher,
		DirectoryFetcher: directoryFetcher,
	}
}

func (s Schema) InputDirectoryToGenerate() error {
	dirPath, err := s.DirectoryFetcher.InteractiveResolveDir(s.Input, constants.MODE_SCHEMA)
	if err != nil {
		return err
	}

	s.DirectoryPath = dirPath
	return nil
}

func (s Schema) InputPropertyNames() error {
	var propertyNames []string
	if err := s.Input.MultipleStringInput(&propertyNames, "Enter property names", s.Validator.Validator_Alphanumeric_Underscore_Allow_Empty()); err != nil {
		return err
	}

	for _, name := range propertyNames {
		property := handler.NewProperty(s.Input, name, s.Property, &handler.Optionals{}, constants.MODE_SCHEMA, s.FileFetcher, s.DirectoryPath)
		s.Properties[name] = property
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

func (s *Schema) GenerateSchema(fileName string, schemaName SchemaName) error {
	data, err := yaml.Marshal(map[SchemaName]*handler.Property{
		schemaName: s.Property,
	})
	if err != nil {
		return err
	}

	if err := utils.GenerateSchema(data, fileName, s.DirectoryPath); err != nil {
		return err
	}

	return nil
}
