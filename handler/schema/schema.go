package schema

import (
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/validator"
)

type SchemaHandler struct {
	Input            input.IInputMethods
	Validator        validator.IInputValidator
	FileFetcher      fetcher.IFileFetcher
	DirectoryFetcher fetcher.IDirectoryFetcher
}

func NewSchemaHandler(input input.IInputMethods, validator validator.IInputValidator, fileFetcher fetcher.IFileFetcher, directoryFetcher fetcher.IDirectoryFetcher) *SchemaHandler {
	return &SchemaHandler{
		Input:            input,
		Validator:        validator,
		FileFetcher:      fileFetcher,
		DirectoryFetcher: directoryFetcher,
	}
}

func (sh *SchemaHandler) HandleGenerateSchemaCommand() error {
	schema := NewSchema(sh.Input, sh.Validator, sh.FileFetcher, sh.DirectoryFetcher)

	if err := schema.InputDirectoryToGenerate(); err != nil {
		return err
	}

	var fileName string
	if err := sh.Input.StringInput(&fileName, "Enter the file name", sh.Validator.Validator_Alphanumeric_Underscore()); err != nil {
		return err
	}

	schema.Type = "object"

	var schemaName SchemaName
	if err := schema.InputSchemaName(&schemaName); err != nil {
		return err
	}

	if err := schema.InputPropertyNames(); err != nil {
		return err
	}

	for _, prop := range schema.Properties {
		if err := prop.ReadAll(); err != nil {
			return err
		}
	}

	if err := schema.GenerateSchema(fileName, schemaName); err != nil {
		return err
	}

	return nil
}
