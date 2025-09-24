package api

import (
	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/handler"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
	"gopkg.in/yaml.v2"
)

type API struct {
	Input              input.IInputMethods       `yaml:"-"`
	APIValidator       validator.IInputValidator `yaml:"-"`
	OptionalProperties handler.Optionals         `yaml:"-"`
	ParameterNames     []string                  `yaml:"-"`
	FileFetcher        fetcher.IFileFetcher      `yaml:"-"`
	DirectoryFetcher   fetcher.IDirectoryFetcher `yaml:"-"`
	DirectoryPath      string                    `yaml:"-"`

	OperationID string               `yaml:"operationId,omitempty"`
	Summary     string               `yaml:"summary,omitempty"`
	Description string               `yaml:"description,omitempty"`
	Tags        []string             `yaml:"tags,omitempty"`
	Parameters  []*Parameter         `yaml:"parameters,omitempty"`
	RequestBody *RequestBody         `yaml:"requestBody,omitempty"`
	Responses   map[string]*Response `yaml:"responses,omitempty"`
}

func NewAPI(input input.IInputMethods, validator validator.IInputValidator, fileFetcher fetcher.IFileFetcher, directoryFetcher fetcher.IDirectoryFetcher) *API {
	return &API{
		Input:            input,
		APIValidator:     validator,
		FileFetcher:      fileFetcher,
		DirectoryFetcher: directoryFetcher,
		Parameters:       []*Parameter{},
		RequestBody:      nil,
		Responses:        make(map[string]*Response),
	}
}

func (a *API) InputDirectoryToGenerate() error {
	dirPath, err := a.DirectoryFetcher.InteractiveResolveDir(a.Input, constants.MODE_API)
	if err != nil {
		return err
	}

	a.DirectoryPath = dirPath
	return nil
}

func (a *API) InputOptionalProperties(method string) error {
	var optionals []string

	if err := a.Input.MultipleSelectInput(&optionals, "Select optional properties", constants.OptionalProperties[method], nil); err != nil {
		return err
	}

	for _, prop := range optionals {
		a.OptionalProperties = append(a.OptionalProperties, prop)
	}

	return nil
}

func (a *API) InputHTTPStatusCodes(method string) error {
	var statusCodes []string
	if err := a.Input.MultipleSelectInput(&statusCodes, "Select HTTP status codes for responses", constants.HTTPStatusMap[method], nil); err != nil {
		return err
	}

	for _, code := range statusCodes {
		a.Responses[code] = NewResponse(a.Input, code, a.OptionalProperties, a.FileFetcher, a.DirectoryPath)
	}

	return nil
}

func (a *API) ReadParameterNames() error {
	var names []string
	if err := a.Input.MultipleStringInput(&names, "Enter parameter names", a.APIValidator.Validator_Alphanumeric_Underscore_Allow_Empty()); err != nil {
		return err
	}
	a.ParameterNames = append(a.ParameterNames, names...)
	return nil
}

func (a *API) ReadOperationID() error {
	if err := a.Input.StringInput(&a.OperationID, "Enter the Operation ID for the API", a.APIValidator.Validator_Alphanumeric_Underscore()); err != nil {
		return err
	}
	return nil
}

func (a *API) ReadSummary() error {
	if err := a.Input.StringInput(&a.Summary, "Enter a brief summary of the API (optional)", nil); err != nil {
		return err
	}
	return nil
}

func (a *API) ReadDescription() error {
	if err := a.Input.StringInput(&a.Description, "Enter a detailed description of the API (optional)", nil); err != nil {
		return err
	}
	return nil
}

func (a *API) ReadTags() error {
	var tagsInput string
	isAdd := true
	for isAdd {
		if err := a.Input.StringInput(&tagsInput, "Enter a tag for the API (optional)", nil); err != nil {
			return err
		}

		if tagsInput != "" {
			a.Tags = append(a.Tags, tagsInput)
		}

		if err := a.Input.BooleanInput(&isAdd, "Do you want to add another tag?"); err != nil {
			return err
		}
	}

	return nil
}

func (a *API) ReadRequestBody() error {
	reqBody := NewRequestBody(a.Input, a.OptionalProperties, a.FileFetcher, a.DirectoryPath)

	if err := reqBody.ReadAll(); err != nil {
		return err
	}

	a.RequestBody = reqBody
	return nil
}

func (a *API) GenerateFile(fileName, method string) error {
	data, err := yaml.Marshal(map[string]*API{
		method: a,
	})
	if err != nil {
		return err
	}

	if err := utils.GenerateSchema(data, fileName, a.DirectoryPath); err != nil {
		return err
	}

	return nil
}

type Parameter struct {
	Input input.IInputMethods `yaml:"-"`

	In     string       `yaml:"in,omitempty"`
	Name   string       `yaml:"name,omitempty"`
	Schema *ParamSchema `yaml:"schema,omitempty"`
}

func NewParameter(input input.IInputMethods, name string, fileFetcher fetcher.IFileFetcher, directoryPath string) *Parameter {
	return &Parameter{
		Input:  input,
		Name:   name,
		Schema: NewParamSchema(input, fileFetcher, directoryPath),
	}
}

func (p *Parameter) ReadIn() error {
	label := "Select Parameter Location (" + p.Name + ")"
	err := p.Input.SelectInput(&p.In, label, constants.ReflableParamIn)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parameter) ReadAll() error {
	if err := p.ReadIn(); err != nil {
		return err
	}

	if err := p.Schema.ReadAll(); err != nil {
		return err
	}

	return nil
}

type ParamSchema struct {
	Input              input.IInputMethods  `yaml:"-"`
	OptionalProperties handler.Optionals    `yaml:"-"`
	FileFetcher        fetcher.IFileFetcher `yaml:"-"`
	DirectoryPath      string               `yaml:"-"`

	Type    string `yaml:"type,omitempty"`
	Format  string `yaml:"format,omitempty"`
	Example string `yaml:"example,omitempty"`
	Max     int64  `yaml:"max,omitempty"`
	Min     int64  `yaml:"min,omitempty"`
	Ref     string `yaml:"$ref,omitempty"`
}

func NewParamSchema(input input.IInputMethods, fileFetcher fetcher.IFileFetcher, directoryPath string) *ParamSchema {
	return &ParamSchema{
		Input:         input,
		FileFetcher:   fileFetcher,
		DirectoryPath: directoryPath,
	}
}

func (ps *ParamSchema) ReadType() error {
	err := ps.Input.SelectInput(&ps.Type, "Select Parameter Type", constants.FieldTypeList)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ParamSchema) ReadFormat() error {
	var format string
	err := ps.Input.SelectInput(&format, "Select Parameter Format", constants.FormatList[ps.Type])
	if err != nil {
		return err
	}

	if ps.Format == constants.FORMAT_NONE {
		return nil
	}

	ps.Format = format
	return nil
}

func (ps *ParamSchema) ReadExample() error {
	var example string
	if err := ps.Input.StringInput(&example, "Enter an example value for the parameter (optional)", nil); err != nil {
		return err
	}

	ps.Example = example
	return nil
}

func (ps *ParamSchema) ReadMax() error {
	var max int64
	if err := ps.Input.Int64Input(&max, "Enter the maximum value for the parameter (optional)", nil); err != nil {
		return err
	}

	ps.Max = max
	return nil
}

func (ps *ParamSchema) ReadMin() error {
	var min int64
	if err := ps.Input.Int64Input(&min, "Enter the minimum value for the parameter (optional)", nil); err != nil {
		return err
	}

	ps.Min = min
	return nil
}

func (ps *ParamSchema) ReadRef() error {
	ref, err := ps.FileFetcher.InteractiveResolveRef(ps.Input, constants.MODE_API, ps.DirectoryPath)
	if err != nil {
		return err
	}

	ps.Ref = ref
	return nil
}

func (ps *ParamSchema) ReadAll() error {
	isReadRef := false
	if err := ps.Input.BooleanInput(&isReadRef, "Do you want to set a $ref for the parameter?"); err != nil {
		return err
	}

	if isReadRef {
		if err := ps.ReadRef(); err != nil {
			return err
		}
		return nil
	}

	if err := ps.ReadType(); err != nil {
		return err
	}

	if constants.IsFormatableType(ps.Type) {
		if err := ps.ReadFormat(); err != nil {
			return err
		}
	}

	if constants.IsExamplableType(ps.Type) && ps.OptionalProperties.Contains(constants.PROPERTY_EXAMPLE) {
		if err := ps.ReadExample(); err != nil {
			return err
		}
	}

	// if constants.IsMaxMinApplicableType(ps.Type) {
	// 	if err := ps.ReadMax(); err != nil {
	// 		return err
	// 	}
	// }

	// if constants.IsMaxMinApplicableType(ps.Type) {
	// 	if err := ps.ReadMin(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type RequestBody struct {
	Input              input.IInputMethods  `yaml:"-"`
	OptionalProperties handler.Optionals    `yaml:"-"`
	FileFetcher        fetcher.IFileFetcher `yaml:"-"`
	DirectoryPath      string               `yaml:"-"`

	Description string                `yaml:"description,omitempty"`
	Required    bool                  `yaml:"required,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
}

func NewRequestBody(input input.IInputMethods, optionalProperties handler.Optionals, fileFetcher fetcher.IFileFetcher, directoryPath string) *RequestBody {
	return &RequestBody{
		Input:              input,
		OptionalProperties: optionalProperties,
		Content:            make(map[string]*MediaType),
		FileFetcher:        fileFetcher,
		DirectoryPath:      directoryPath,
	}
}

func (rq *RequestBody) InputMediaTypes() error {
	var mediaTypes []string
	label := "Select media types for the request body"
	if err := rq.Input.MultipleSelectInput(&mediaTypes, label, constants.MimeKeys, nil); err != nil {
		return err
	}

	for _, mt := range mediaTypes {
		mimeType := constants.MediaTypeMap[mt]
		rq.Content[mimeType] = NewMediaType(rq.Input, rq.OptionalProperties, rq.FileFetcher, rq.DirectoryPath)
	}

	return nil
}

func (rq *RequestBody) ReadDescription() error {
	label := "Enter a description for the request body"
	if err := rq.Input.StringInput(&rq.Description, label, nil); err != nil {
		return err
	}
	return nil
}

func (rq *RequestBody) ReadRequired() error {
	if err := rq.Input.BooleanInput(&rq.Required, "Is the request body required?"); err != nil {
		return err
	}
	return nil
}

func (rq *RequestBody) ReadAll() error {
	if rq.OptionalProperties.Contains(constants.PROPERTY_DESCRIPTION) {
		if err := rq.ReadDescription(); err != nil {
			return err
		}
	}

	if err := rq.ReadRequired(); err != nil {
		return err
	}

	if err := rq.InputMediaTypes(); err != nil {
		return err
	}

	for _, mt := range rq.Content {
		if err := mt.ReadAll(); err != nil {
			return err
		}
	}
	return nil
}

type MediaType struct {
	Input input.IInputMethods `yaml:"-"`

	Schema *handler.Property `yaml:"schema,omitempty"`
	// TODO: Implement examples
	// Example string `yaml:"example,omitempty"`
}

func NewMediaType(input input.IInputMethods, optionalProperties handler.Optionals, fileFetcher fetcher.IFileFetcher, directoryPath string) *MediaType {
	return &MediaType{
		Input:  input,
		Schema: handler.NewProperty(input, "schema", nil, &optionalProperties, constants.MODE_API, fileFetcher, directoryPath),
	}
}

func (mt *MediaType) ReadAll() error {
	if err := mt.Schema.ReadAll(); err != nil {
		return err
	}

	return nil
}

type Response struct {
	Input              input.IInputMethods  `yaml:"-"`
	Code               string               `yaml:"-"`
	OptionalProperties handler.Optionals    `yaml:"-"`
	FileFetcher        fetcher.IFileFetcher `yaml:"-"`
	DirectoryPath      string               `yaml:"-"`

	Description string                `yaml:"description,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
}

func NewResponse(input input.IInputMethods, code string, optionalProperties handler.Optionals, fileFetcher fetcher.IFileFetcher, directoryPath string) *Response {
	return &Response{
		Input:              input,
		Code:               code,
		OptionalProperties: optionalProperties,
		Content:            make(map[string]*MediaType),
		FileFetcher:        fileFetcher,
		DirectoryPath:      directoryPath,
	}
}

func (r *Response) InputMediaTypes() error {
	var mediaTypes []string
	label := "Select media types for the response (" + r.Code + ")"
	if err := r.Input.MultipleSelectInput(&mediaTypes, label, constants.MimeKeys, nil); err != nil {
		return err
	}

	for _, mt := range mediaTypes {
		mimeType := constants.MediaTypeMap[mt]
		r.Content[mimeType] = NewMediaType(r.Input, r.OptionalProperties, r.FileFetcher, r.DirectoryPath)
	}

	return nil
}

func (r *Response) ReadDescription(code string) error {
	label := "Enter a description for the response (" + code + ")"
	if err := r.Input.StringInput(&r.Description, label, nil); err != nil {
		return err
	}
	return nil
}

func (r *Response) ReadAll(code string, isReadDescription bool) error {
	if isReadDescription {
		if err := r.ReadDescription(code); err != nil {
			return err
		}
	}

	if err := r.InputMediaTypes(); err != nil {
		return err
	}

	for _, mt := range r.Content {
		if err := mt.ReadAll(); err != nil {
			return err
		}
	}
	return nil
}
