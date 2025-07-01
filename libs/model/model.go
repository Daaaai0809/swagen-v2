package model

import (
	"github.com/Daaaai0809/swagen-v2/libs"
	"github.com/Daaaai0809/swagen-v2/utils"
)

type ModelHandler struct {
	Input utils.IInputMethods
}

func NewModelHandler(input utils.IInputMethods) *ModelHandler {
	return &ModelHandler{
		Input: input,
	}
}

func (mh *ModelHandler) HandleGenerateModelCommand() error {
	model := NewModel(mh.Input)

	if err := model.ReadTitle(); err != nil {
		return err
	}

	for {
		var propertyName string
		if err := model.ReadPropertyName(&propertyName); err != nil {
			return err
		}

		schema := libs.NewSchema(mh.Input, propertyName)
		if err := schema.ReadSchema(); err != nil {
			return err
		}
		model.Properties[propertyName] = schema

		isAdd := false
		if err := mh.Input.BooleanInput(&isAdd, "Do you want to add another property? (Root Property)"); err != nil {
			return err
		}
		if !isAdd {
			break
		}
	}

	if err := utils.GenerateSchema(model, model.Title, utils.GetEnv("MODEL_PATH", "models/")); err != nil {
		return err
	}

	return nil
}
