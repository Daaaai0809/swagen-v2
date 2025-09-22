package schema

import (
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
)

type SchemaHandler struct {
	Input       input.IInputMethods
	Validator   validator.IInputValidator
	FileFetcher fetcher.IFileFetcher
}

func NewSchemaHandler(input input.IInputMethods, validator validator.IInputValidator, fileFetcher fetcher.IFileFetcher) *SchemaHandler {
	return &SchemaHandler{
		Input:       input,
		Validator:   validator,
		FileFetcher: fileFetcher,
	}
}

func (sh *SchemaHandler) HandleGenerateSchemaCommand() error {
	var fileName string
	if err := sh.Input.StringInput(&fileName, "Enter the file name", sh.Validator.Validator_Alphanumeric_Underscore()); err != nil {
		return err
	}

	schema := NewSchema(sh.Input, sh.Validator, sh.FileFetcher)

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

	if err := schema.GenerateSchema(fileName, schemaName, utils.GetEnv(utils.SWAGEN_SCHEMA_PATH, "schema/")); err != nil {
		return err
	}

	return nil
}
