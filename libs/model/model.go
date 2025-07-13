package model

import (
	"errors"
	"strings"

	"github.com/Daaaai0809/swagen-v2/constants"
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

	var validate utils.ValidationFunc = func(input string) error {
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
	if err := mh.Input.StringInput(&fileName, "Enter the model file name (without extension)", &validate); err != nil {
		return err
	}

	if err := model.ReadTitle(); err != nil {
		return err
	}

	for {
		var propertyName string
		if err := model.ReadPropertyName(&propertyName); err != nil {
			return err
		}

		property := libs.NewProperty(mh.Input, propertyName, nil, constants.MODE_MODEL)
		if err := property.ReadAll(); err != nil {
			return err
		}
		model.Properties[propertyName] = property

		isAdd := false
		if err := mh.Input.BooleanInput(&isAdd, "Do you want to add another property? (Root Property)"); err != nil {
			return err
		}
		if !isAdd {
			break
		}
	}

	if err := model.GenerateModel(fileName, utils.GetEnv("MODEL_PATH", "models/")); err != nil {
		return err
	}

	return nil
}
