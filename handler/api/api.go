package api

import (
	"os"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
	"gopkg.in/yaml.v2"
)

type APIMap map[string]*API

func (am *APIMap) GetMethods() []string {
	methods := make([]string, 0, len(*am))
	for method := range *am {
		methods = append(methods, method)
	}
	return methods
}

func (am *APIMap) ToYaml() ([]byte, error) {
	return yaml.Marshal(am)
}

type IAPIHandler interface {
	HandleGenerateAPICommand() error
	HandleAddToAPICommand() error
}

type APIHandler struct {
	Input            input.IInputMethods
	APIValidator     validator.IInputValidator
	FileFetcher      fetcher.IFileFetcher
	DirectoryFetcher fetcher.IDirectoryFetcher
}

func NewAPIHandler(input input.IInputMethods, validator validator.IInputValidator, fileFetcher fetcher.IFileFetcher, directoryFetcher fetcher.IDirectoryFetcher) IAPIHandler {
	return &APIHandler{
		Input:            input,
		APIValidator:     validator,
		FileFetcher:      fileFetcher,
		DirectoryFetcher: directoryFetcher,
	}
}

func (ah *APIHandler) HandleGenerateAPICommand() error {
	api := NewAPI(ah.Input, ah.APIValidator, ah.FileFetcher, ah.DirectoryFetcher)

	if err := api.InputDirectoryToGenerate(); err != nil {
		return err
	}

	var fileName string
	if err := ah.Input.StringInput(&fileName, "Enter the API file name (without extension)", ah.APIValidator.Validator_Alphanumeric_Underscore()); err != nil {
		return err
	}

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
			param := NewParameter(api.Input, name, api.FileFetcher, api.DirectoryPath)
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

	if err := api.GenerateFile(fileName, method); err != nil {
		return err
	}

	return nil
}

func (ah *APIHandler) HandleAddToAPICommand() error {
	filePath, directoryPath, err := ah.FileFetcher.FetchPathSchema(ah.Input)
	if err != nil {
		return err
	}

	existingAPI, err := ah.parseExistingAPI(filePath)
	if err != nil {
		return err
	}

	api := NewAPI(ah.Input, ah.APIValidator, ah.FileFetcher, ah.DirectoryFetcher)

	api.DirectoryPath = directoryPath

	var method string
	if err := ah.Input.SelectInput(&method, "Select the HTTP method for the API", constants.GetNotExistingMethods(existingAPI.GetMethods())); err != nil {
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
			param := NewParameter(api.Input, name, api.FileFetcher, api.DirectoryPath)
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

	existingAPI[constants.HTTPMethodsMap[method]] = api

	yamlData, err := existingAPI.ToYaml()
	if err != nil {
		return err
	}

	if err := utils.WriteToFile(yamlData, filePath); err != nil {
		return err
	}
	return nil
}

// parseExistingAPI parses an existing API schema file and returns the API object map
func (ah *APIHandler) parseExistingAPI(filePath string) (APIMap, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var apiMap APIMap
	if err := yaml.Unmarshal(data, &apiMap); err != nil {
		return nil, err
	}

	return apiMap, nil
}
