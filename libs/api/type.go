package api

import (
	"errors"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/utils"
	"gopkg.in/yaml.v2"
)

type API struct {
	Input utils.IInputMethods `yaml:"-"`

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
	Input utils.IInputMethods `yaml:"-"`

	In     string       `yaml:"in,omitempty"`
	Name   string       `yaml:"name,omitempty"`
	Schema *ParamSchema `yaml:"schema,omitempty"`
}

func NewParameter(input utils.IInputMethods) *Parameter {
	return &Parameter{
		Input:  input,
		Schema: newParamSchema(input),
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
	var validate utils.ValidationFunc = func(input string) error {
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
	Input utils.IInputMethods `yaml:"-"`

	Type    string `yaml:"type,omitempty"`
	Format  string `yaml:"format,omitempty"`
	Example string `yaml:"example,omitempty"`
	Max     int64  `yaml:"max,omitempty"`
	Min     int64  `yaml:"min,omitempty"`
	Ref     string `yaml:"$ref,omitempty"`
}

func newParamSchema(input utils.IInputMethods) *ParamSchema {
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
	Input utils.IInputMethods `yaml:"-"`

	Description string                `yaml:"description,omitempty"`
	Required    bool                  `yaml:"required,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
}

func NewRequestBody(input utils.IInputMethods) *RequestBody {
	return &RequestBody{
		Input:   input,
		Content: make(map[string]*MediaType),
	}
}

type MediaType struct {
	Input utils.IInputMethods `yaml:"-"`

	Schema *MediaTypeSchema `yaml:"schema,omitempty"`
	// TODO: Implement examples
	// Example string `yaml:"example,omitempty"`
}

func NewMediaType(input utils.IInputMethods) *MediaType {
	return &MediaType{
		Input:  input,
		Schema: newMediaTypeSchema(input),
	}
}

func (mt *MediaType) ReadAll() error {
	if err := mt.Schema.ReadAll(); err != nil {
		return err
	}

	return nil
}

type MediaTypeSchema struct {
	Input utils.IInputMethods `yaml:"-"`

	Type   string `yaml:"type,omitempty"`
	Format string `yaml:"format,omitempty"`
	Ref    string `yaml:"$ref,omitempty"`
}

func newMediaTypeSchema(input utils.IInputMethods) *MediaTypeSchema {
	return &MediaTypeSchema{
		Input: input,
	}
}

func (mts *MediaTypeSchema) ReadType() error {
	err := mts.Input.SelectInput(&mts.Type, "Select Schema Type", constants.FieldTypeList)
	if err != nil {
		return err
	}

	return nil
}

func (mts *MediaTypeSchema) ReadFormat() error {
	err := mts.Input.SelectInput(&mts.Format, "Select Schema Format", constants.FormatList[mts.Type])
	if err != nil {
		return err
	}

	return nil
}

func (mts *MediaTypeSchema) ReadRef() error {
	ref, err := utils.InteractiveResolveRef(mts.Input, constants.MODE_API)
	if err != nil {
		return err
	}

	mts.Ref = ref
	return nil
}

func (mts *MediaTypeSchema) ReadAll() error {
	isReadRef := false
	if err := mts.Input.BooleanInput(&isReadRef, "Do you want to set a $ref for the schema?"); err != nil {
		return err
	}
	if isReadRef {
		if err := mts.ReadRef(); err != nil {
			return err
		}
		return nil
	}

	if err := mts.ReadType(); err != nil {
		return err
	}

	isReadFormat := false
	if err := mts.Input.BooleanInput(&isReadFormat, "Do you want to set a format for the schema?"); err != nil {
		return err
	}
	if isReadFormat && constants.IsFormatableType(mts.Type) {
		if err := mts.ReadFormat(); err != nil {
			return err
		}
	}

	return nil
}

type Response struct {
	Input utils.IInputMethods `yaml:"-"`

	Description string                `yaml:"description,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
}

func NewResponse(input utils.IInputMethods) *Response {
	return &Response{
		Input:   input,
		Content: make(map[string]*MediaType),
	}
}
