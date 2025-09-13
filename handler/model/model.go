package model

import (
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
)

type ModelHandler struct {
	Input     input.IInputMethods
	Validator validator.IInputValidator
}

func NewModelHandler(input input.IInputMethods, validator validator.IInputValidator) *ModelHandler {
	return &ModelHandler{
		Input:     input,
		Validator: validator,
	}
}

func (mh *ModelHandler) HandleGenerateModelCommand() error {
	model := NewModel(mh.Input, mh.Validator)

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

	if err := model.GenerateModel(fileName, utils.GetEnv(utils.SWAGEN_MODEL_PATH, "models/")); err != nil {
		return err
	}

	return nil
}
