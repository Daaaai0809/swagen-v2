package model

import (
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/validator"
)

type ModelHandler struct {
	Input            input.IInputMethods
	Validator        validator.IInputValidator
	DirectoryFetcher fetcher.IDirectoryFetcher
}

func NewModelHandler(input input.IInputMethods, validator validator.IInputValidator, directoryFetcher fetcher.IDirectoryFetcher) *ModelHandler {
	return &ModelHandler{
		Input:            input,
		Validator:        validator,
		DirectoryFetcher: directoryFetcher,
	}
}

func (mh *ModelHandler) HandleGenerateModelCommand() error {
	model := NewModel(mh.Input, mh.Validator, mh.DirectoryFetcher)

	if err := model.InputDirectoryToGenerate(); err != nil {
		return err
	}

	var fileName string
	if err := mh.Input.StringInput(&fileName, "Enter the model file name (without extension)", mh.Validator.Validator_Alphanumeric_Underscore()); err != nil {
		return err
	}

	if err := model.ReadTitle(); err != nil {
		return err
	}

	if err := model.ReadPropertyNames(); err != nil {
		return err
	}

	for _, prop := range model.Properties {
		if err := prop.ReadAll(); err != nil {
			return err
		}
	}

	if err := model.GenerateModel(fileName); err != nil {
		return err
	}

	return nil
}
