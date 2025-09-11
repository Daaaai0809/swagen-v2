package api

import (
	"errors"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/handler"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"gopkg.in/yaml.v2"
)

type API struct {
	Input input.IInputMethods `yaml:"-"`

	OperationID string               `yaml:"operationId,omitempty"`
	Summary     string               `yaml:"summary,omitempty"`
	Description string               `yaml:"description,omitempty"`
	Tags        []string             `yaml:"tags,omitempty"`
	Parameters  []*Parameter         `yaml:"parameters,omitempty"`
	RequestBody *RequestBody         `yaml:"requestBody,omitempty"`
	Responses   map[string]*Response `yaml:"responses,omitempty"`
}

func NewAPI() *API {
	return &API{
		Parameters:  []*Parameter{},
		RequestBody: nil,
		Responses:   make(map[string]*Response),
	}
}

func (a *API) ReadOperationID() error {
	if err := a.Input.StringInput(&a.OperationID, "Enter the Operation ID for the API (optional)", nil); err != nil {
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

func (a *API) ReadParameter() error {
	param := NewParameter(a.Input)

	if err := param.ReadIn(); err != nil {
		return err
	}

	if err := param.ReadName(); err != nil {
		return err
	}

	if err := param.Schema.ReadAll(); err != nil {
		return err
	}

	a.Parameters = append(a.Parameters, param)

	return nil
}

func (a *API) ReadRequestBody() error {
	reqBody := NewRequestBody(a.Input)

	var description string
	if err := reqBody.Input.StringInput(&description, "Enter a description for the request body (optional)", nil); err != nil {
		return err
	}
	reqBody.Description = description

	if err := reqBody.Input.BooleanInput(&reqBody.Required, "Is the request body required?"); err != nil {
		return err
	}

	for {
		var mediaType string
		if err := reqBody.Input.SelectInput(&mediaType, "Select a media type for the request body (or leave empty to finish)", constants.MediaTypeList); err != nil {
			return err
		}

		if mediaType == "" {
			break
		}

		mt := NewMediaType(a.Input)
		if err := mt.Schema.ReadAll(); err != nil {
			return err
		}

		reqBody.Content[mediaType] = mt

		var addMore bool
		if err := a.Input.BooleanInput(&addMore, "Do you want to add another media type?"); err != nil {
			return err
		}
		if !addMore {
			break
		}
	}

	a.RequestBody = reqBody
	return nil
}

func (a *API) ReadResponse() error {
	resp := NewResponse(a.Input)

	var httpStatus string
	if err := resp.Input.SelectInput(&httpStatus, "Select the HTTP status code for the response", constants.HTTPStatusList); err != nil {
		return err
	}

	var description string
	if err := resp.Input.StringInput(&description, "Enter a description for the response (optional)", nil); err != nil {
		return err
	}
	resp.Description = description

	for {
		var mediaType string
		if err := resp.Input.SelectInput(&mediaType, "Select a media type for the response (or leave empty to finish)", constants.MediaTypeList); err != nil {
			return err
		}

		if mediaType == "" {
			break
		}

		mt := NewMediaType(a.Input)
		if err := mt.Schema.ReadAll(); err != nil {
			return err
		}

		resp.Content[mediaType] = mt

		var addMore bool
		if err := a.Input.BooleanInput(&addMore, "Do you want to add another media type?"); err != nil {
			return err
		}
		if !addMore {
			break
		}
	}

	a.Responses[httpStatus] = resp
	return nil
}

func (a *API) GenerateFile(fileName, method, path string) error {
	data, err := yaml.Marshal(map[string]*API{
		method: a,
	})
	if err != nil {
		return err
	}

	if err := utils.GenerateSchema(data, fileName, path); err != nil {
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

func NewParameter(input input.IInputMethods) *Parameter {
	return &Parameter{
		Input:  input,
		Schema: NewParamSchema(input),
	}
}

func (p *Parameter) ReadIn() error {
	err := p.Input.SelectInput(&p.In, "Select Parameter Location", constants.ReflableParamIn)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parameter) ReadName() error {
	var validate input.ValidationFunc = func(input string) error {
		if input == "" {
			return errors.New("parameter name is required")
		}
		return nil
	}

	if err := p.Input.StringInput(&p.Name, "Enter the parameter name", &validate); err != nil {
		return err
	}

	return nil
}

type ParamSchema struct {
	Input input.IInputMethods `yaml:"-"`

	Type    string `yaml:"type,omitempty"`
	Format  string `yaml:"format,omitempty"`
	Example string `yaml:"example,omitempty"`
	Max     int64  `yaml:"max,omitempty"`
	Min     int64  `yaml:"min,omitempty"`
	Ref     string `yaml:"$ref,omitempty"`
}

func NewParamSchema(input input.IInputMethods) *ParamSchema {
	return &ParamSchema{
		Input: input,
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
	err := ps.Input.SelectInput(&ps.Format, "Select Parameter Format", constants.FormatList[ps.Type])
	if err != nil {
		return err
	}

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
	ref, err := utils.InteractiveResolveRef(ps.Input, constants.MODE_API)
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

	isReadFormat := false
	if err := ps.Input.BooleanInput(&isReadFormat, "Do you want to set a format for the parameter?"); err != nil {
		return err
	}
	if isReadFormat && constants.IsFormatableType(ps.Type) {
		if err := ps.ReadFormat(); err != nil {
			return err
		}
	}

	isReadExample := false
	if err := ps.Input.BooleanInput(&isReadExample, "Do you want to add an example value for the parameter?"); err != nil {
		return err
	}
	if isReadExample && constants.IsExamplableType(ps.Type) {
		if err := ps.ReadExample(); err != nil {
			return err
		}
	}

	isReadMax := false
	if err := ps.Input.BooleanInput(&isReadMax, "Do you want to set a maximum value for the parameter?"); err != nil {
		return err
	}
	if isReadMax && constants.IsMaxMinApplicableType(ps.Type) {
		if err := ps.ReadMax(); err != nil {
			return err
		}
	}

	isReadMin := false
	if err := ps.Input.BooleanInput(&isReadMin, "Do you want to set a minimum value for the parameter?"); err != nil {
		return err
	}
	if isReadMin && constants.IsMaxMinApplicableType(ps.Type) {
		if err := ps.ReadMin(); err != nil {
			return err
		}
	}

	return nil
}

type RequestBody struct {
	Input input.IInputMethods `yaml:"-"`

	Description string                `yaml:"description,omitempty"`
	Required    bool                  `yaml:"required,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
}

func NewRequestBody(input input.IInputMethods) *RequestBody {
	return &RequestBody{
		Input:   input,
		Content: make(map[string]*MediaType),
	}
}

type MediaType struct {
	Input input.IInputMethods `yaml:"-"`

	Schema *handler.Property `yaml:"schema,omitempty"`
	// TODO: Implement examples
	// Example string `yaml:"example,omitempty"`
}

func NewMediaType(input input.IInputMethods) *MediaType {
	return &MediaType{
		Input:  input,
		Schema: handler.NewProperty(input, "schema", nil, constants.MODE_API),
	}
}

func (mt *MediaType) ReadAll() error {
	if err := mt.Schema.ReadAll(); err != nil {
		return err
	}

	return nil
}

type Response struct {
	Input input.IInputMethods `yaml:"-"`

	Description string                `yaml:"description,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
}

func NewResponse(input input.IInputMethods) *Response {
	return &Response{
		Input:   input,
		Content: make(map[string]*MediaType),
	}
}
