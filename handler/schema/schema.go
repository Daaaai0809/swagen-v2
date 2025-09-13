package schema

import (
	"errors"
	"strings"

	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
)

type SchemaHandler struct {
	Input input.IInputMethods
}

func NewSchemaHandler(input input.IInputMethods) *SchemaHandler {
	return &SchemaHandler{
		Input: input,
	}
}

func (sh *SchemaHandler) HandleGenerateSchemaCommand() error {
	var validate input.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("file name is required")
		}

		// NOTE: 数字スタートを許可しない
		if strings.HasPrefix(input, "0") {
			return errors.New("file name cannot start with a number")
		}

		// NOTE: 英数字とアンダースコアのみを許可
		for _, char := range input {
			if !(('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || char == '_') {
				return errors.New("file name can only contain alphanumeric characters and underscores")
			}
		}

		return nil
	}

	var fileName string
	if err := sh.Input.StringInput(&fileName, ">", &validate); err != nil {
		return err
	}

	schema := NewSchema(sh.Input)

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

	if err := schema.GenerateSchema(fileName, schemaName, utils.GetEnv(utils.SCHEMA_PATH, "schema/")); err != nil {
		return err
	}

	return nil
}
