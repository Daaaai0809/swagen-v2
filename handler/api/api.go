package api

import (
	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
)

type APIHandler struct {
	Input        input.IInputMethods
	APIValidator validator.IInputValidator
}

func NewAPIHandler(input input.IInputMethods, validator validator.IInputValidator) *APIHandler {
	return &APIHandler{
		Input:        input,
		APIValidator: validator,
	}
}

func (ah *APIHandler) HandleGenerateAPICommand() error {
	var fileName string
	if err := ah.Input.StringInput(&fileName, "Enter the API file name (without extension)", ah.APIValidator.Validator_Alphanumeric_Underscore()); err != nil {
		return err
	}

	api := NewAPI(ah.Input, ah.APIValidator)

	var method string
	if err := ah.Input.SelectInput(&method, "Select the HTTP method for the API", constants.HTTPMethods); err != nil {
		return err
	}

	if err := api.InputOptionalProperties(method); err != nil {
		return err
	}

	if api.OptionalProperties.Contains(constants.PROPERTY_OPERATION_ID) {
		if err := api.ReadOperationID(); err != nil {
			return err
		}
	}

	if api.OptionalProperties.Contains(constants.PROPERTY_SUMMARY) {
		if err := api.ReadSummary(); err != nil {
			return err
		}
	}

	if api.OptionalProperties.Contains(constants.PROPERTY_DESCRIPTION) {
		if err := api.ReadDescription(); err != nil {
			return err
		}
	}

	if api.OptionalProperties.Contains(constants.PROPERTY_TAGS) {
		if err := api.ReadTags(); err != nil {
			return err
		}
	}

	if api.OptionalProperties.Contains(constants.PROPERTY_PARAMETERS) {
		if err := api.ReadParameterNames(); err != nil {
			return err
		}

		for _, name := range api.ParameterNames {
			param := NewParameter(api.Input, name)
			api.Parameters = append(api.Parameters, param)
			if err := param.ReadAll(); err != nil {
				return err
			}
		}
	}

	if api.OptionalProperties.Contains(constants.PROPERTY_REQUEST_BODY) {
		if err := api.ReadRequestBody(); err != nil {
			return err
		}
	}

	if err := api.InputHTTPStatusCodes(method); err != nil {
		return err
	}

	for code, resp := range api.Responses {
		if err := resp.ReadAll(code, api.OptionalProperties.Contains(constants.PROPERTY_DESCRIPTION)); err != nil {
			return err
		}
	}

	if err := api.GenerateFile(fileName, constants.HTTPMethodsMap[method], utils.GetEnv(utils.SWAGEN_API_PATH, "path/")); err != nil {
		return err
	}

	return nil
}
