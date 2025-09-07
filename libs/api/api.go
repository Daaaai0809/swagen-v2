package api

import (
	"errors"
	"strings"

	"github.com/Daaaai0809/swagen-v2/utils"
)

type APIHandler struct {
	Input utils.IInputMethods
}

func NewAPIHandler(input utils.IInputMethods) *APIHandler {
	return &APIHandler{
		Input: input,
	}
}

func (ah *APIHandler) HandleGenerateAPICommand() error {
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
	if err := ah.Input.StringInput(&fileName, "Enter the API file name (without extension)", &validate); err != nil {
		return err
	}

	api := NewAPI()
	
	if err := api.ReadOperationID(); err != nil {
		return err
	}

	var method string
	if err := ah.Input.SelectInput(&method, "Select the HTTP method for the API", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}); err != nil {
		return err
	}

	isReadSummary := false
	if err := ah.Input.BooleanInput(&isReadSummary, "Do you want to add a summary?"); err != nil {
		return err
	}
	if isReadSummary {
		if err := api.ReadSummary(); err != nil {
			return err
		}
	}

	isReadDescription := false
	if err := ah.Input.BooleanInput(&isReadDescription, "Do you want to add a description?"); err != nil {
		return err
	}
	if isReadDescription {
		if err := api.ReadDescription(); err != nil {
			return err
		}
	}
	
	isReadTags := false
	if err := ah.Input.BooleanInput(&isReadTags, "Do you want to add tags?"); err != nil {
		return err
	}
	if isReadTags {
		if err := api.ReadTags(); err != nil {
			return err
		}
	}

	for {
		var addParam bool
		if err := ah.Input.BooleanInput(&addParam, "Do you want to add a parameter?"); err != nil {
			return err
		}

		if !addParam {
			break
		}

		if err := api.ReadParameter(); err != nil {
			return err
		}
	}

	isReadRequestBody := false
	if err := ah.Input.BooleanInput(&isReadRequestBody, "Do you want to add a request body?"); err != nil {
		return err
	}
	if isReadRequestBody {
		if err := api.ReadRequestBody(); err != nil {
			return err
		}
	}

	for {
		var addResponse bool
		if err := ah.Input.BooleanInput(&addResponse, "Do you want to add a response?"); err != nil {
			return err
		}

		if !addResponse {
			break
		}

		if err := api.ReadResponse(); err != nil {
			return err
		}
	}

	if err := api.GenerateFile(fileName, method, utils.GetEnv(utils.API_PATH, "path/")); err != nil {
		return err
	}

	return nil
}
