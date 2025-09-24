package model

import (
	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/handler"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
	"gopkg.in/yaml.v2"
)

type Model struct {
	Input            input.IInputMethods       `yaml:"-"`
	Validator        validator.IInputValidator `yaml:"-"`
	DirectoryFetcher fetcher.IDirectoryFetcher `yaml:"-"`
	DirectoryPath    string                    `yaml:"-"`

	Title      string                       `yaml:"title,omitempty"`
	Type       string                       `yaml:"type"`
	Properties map[string]*handler.Property `yaml:"properties,omitempty"`
}

func NewModel(input input.IInputMethods, validator validator.IInputValidator, directoryFetcher fetcher.IDirectoryFetcher) *Model {
	return &Model{
		Input:            input,
		Validator:        validator,
		DirectoryFetcher: directoryFetcher,
		Title:            "",
		Type:             constants.OBJECT_TYPE,
		Properties:       make(map[string]*handler.Property),
	}
}

func (m *Model) InputDirectoryToGenerate() error {
	dirPath, err := m.DirectoryFetcher.InteractiveResolveDir(m.Input, constants.MODE_MODEL)
	if err != nil {
		return err
	}

	m.DirectoryPath = dirPath
	return nil
}

func (m *Model) ReadTitle() error {
	err := m.Input.StringInput(&m.Title, "Enter the model title", nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) ReadPropertyNames() error {
	var propertyNames []string
	if err := m.Input.MultipleStringInput(&propertyNames, "Enter property names", m.Validator.Validator_Alphanumeric_Underscore_Allow_Empty()); err != nil {
		return err
	}

	for _, name := range propertyNames {
		property := handler.NewProperty(m.Input, name, nil, &handler.Optionals{}, constants.MODE_MODEL, nil)
		m.Properties[name] = property
	}

	return nil
}

func (m *Model) GenerateModel(fileName string) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	if err := utils.GenerateSchema(data, fileName, m.DirectoryPath); err != nil {
		return err
	}

	return nil
}
