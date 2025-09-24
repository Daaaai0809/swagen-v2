package fetcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"gopkg.in/yaml.v2"
)

// FileFetcher specific constants
const (
	DEFAULT_DEST_PATH       = "."
	PROPERTIES_PATH         = "/properties"
	ITEMS_PATH              = "/items"
	CONTINUE_INTO_SUB_PROPS = "This is an object. Continue into sub-properties?"
	ARRAY_ITEMS_OR_USE      = "This is an array. Select items or use array as is?"
	ARRAY_DETECTED_MSG      = "Array detected. Select items or use array as is?"
	USE_ROOT_SCHEMA_MSG     = "Use the root schema as reference?"
	SELECT_PROPERTY_MSG     = "Select property"
	SELECT_ROOT_SCHEMA_MSG  = "Select root schema"
	WHICH_BASE_PATH_MSG     = "Which base path to reference?"
	MODEL                   = "MODEL"
	SCHEMA                  = "SCHEMA"
)

type IFileFetcher interface {
	InteractiveResolveRef(input input.IInputMethods, mode constants.InputMode, destBase string) (string, error)
}

// FileFetcher handles file-specific fetching operations
type FileFetcher struct {
	baseFetcher IBaseFetcher
}

// NewFileFetcher creates a new FileFetcher instance
func NewFileFetcher() *FileFetcher {
	return &FileFetcher{
		baseFetcher: NewBaseFetcher(),
	}
}

// InteractiveResolveRef
// Contract:
// - Inputs: input.IInputMethods, mode (SCHEMA or API)
// - Outcome: returns a $ref string built from a user-selected YAML file and field
// - Behavior:
//   - SCHEMA mode: traverse from SWAGEN_MODEL_PATH
//   - API mode: ask user to choose start path (SWAGEN_MODEL_PATH or SWAGEN_SCHEMA_PATH)
//   - After selecting a file, parse YAML into Model or Schema(map) and select a field
//   - Build JSON Pointer and relative file path (relative to SCHEMA_PATH or API_PATH)
func (ff *FileFetcher) InteractiveResolveRef(input input.IInputMethods, mode constants.InputMode, destBase string) (string, error) {
	// Decide start path and destination base to compute relative path
	startPath, fileKind, err := ff.decideStartPath(input, mode)
	if err != nil {
		return "", err
	}

	// Directory traversal to pick a YAML file
	selectedFile, err := ff.selectYamlFileInteractive(input, startPath)
	if err != nil {
		return "", err
	}

	// If startPath came from SWAGEN_SCHEMA_PATH, treat file as schema-kind even in API mode
	if fileKind == "auto" {
		// infer by extension only (already .yaml) and location
		if strings.HasPrefix(filepath.Clean(selectedFile), filepath.Clean(utils.GetEnv(utils.SWAGEN_SCHEMA_PATH, ""))) {
			fileKind = "schema"
		} else {
			fileKind = "model"
		}
	}

	// Parse YAML and select a field to build JSON Pointer
	var pointer string
	switch fileKind {
	case "model":
		pointer, err = ff.selectFieldFromModelFile(input, selectedFile)
	case "schema":
		pointer, err = ff.selectFieldFromSchemaFile(input, selectedFile)
	default:
		err = errors.New("unknown file kind")
	}
	if err != nil {
		return "", err
	}

	// Build relative path
	rel, err := filepath.Rel(destBase, selectedFile)
	if err != nil {
		return "", fmt.Errorf("[ERROR] relative path resolution failed")
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasSuffix(rel, YAML_EXT) && !strings.HasSuffix(rel, YML_EXT) {
		// safety: ensure extension
		rel += YAML_EXT
	}

	// Ensure pointer starts with '#'
	if pointer == "" {
		pointer = JSON_POINTER_REF
	} else if !strings.HasPrefix(pointer, JSON_POINTER_REF) {
		pointer = JSON_POINTER_REF + pointer
	}

	return fmt.Sprintf("%s%s", rel, pointer), nil
}

// decideStartPath asks for start directory based on mode and returns also the fileKind hint
// fileKind: "model", "schema", or "auto"
func (ff *FileFetcher) decideStartPath(input input.IInputMethods, mode constants.InputMode) (string, string, error) {
	switch mode {
	case constants.MODE_SCHEMA:
		p := utils.GetEnv(utils.SWAGEN_MODEL_PATH, "")
		if p == "" {
			return "", "", errors.New("[ERROR] SWAGEN_MODEL_PATH is not set. Set it in environment or .env")
		}
		return p, "model", nil
	case constants.MODE_API:
		var choice string
		if err := input.SelectInput(&choice, WHICH_BASE_PATH_MSG, []string{MODEL, SCHEMA}); err != nil {
			return "", "", err
		}
		switch choice {
		case MODEL:
			p := utils.GetEnv(utils.SWAGEN_MODEL_PATH, "")
			if p == "" {
				return "", "", errors.New("[ERROR] SWAGEN_MODEL_PATH is not set. Set it in environment or .env")
			}
			return p, "model", nil
		case SCHEMA:
			p := utils.GetEnv(utils.SWAGEN_SCHEMA_PATH, "")
			if p == "" {
				return "", "", errors.New("[ERROR] SWAGEN_SCHEMA_PATH is not set. Set it in environment or .env")
			}
			return p, "schema", nil
		default:
			return "", "", errors.New("[ERROR] invalid base path selection")
		}
	default:
		return "", "", errors.New("[ERROR] unsupported mode for ref resolution")
	}
}

// selectYamlFileInteractive lets the user navigate directories and select a YAML file.
func (ff *FileFetcher) selectYamlFileInteractive(input input.IInputMethods, start string) (string, error) {
	yamlFilter := func(filename string) bool {
		lower := strings.ToLower(filename)
		return strings.HasSuffix(lower, YAML_EXT) || strings.HasSuffix(lower, YML_EXT)
	}

	return ff.baseFetcher.SelectFileInteractive(input, start, yamlFilter)
}

// selectFieldFromModelFile parses a model file and guides the user to select a field.
// Returns a JSON Pointer like "/properties/foo/items/properties/bar" (without leading '#').
// local lite types to avoid importing handler and causing cycles
type propertyLite struct {
	Type       string                   `yaml:"type,omitempty"`
	Properties map[string]*propertyLite `yaml:"properties,omitempty"`
	Items      *propertyLite            `yaml:"items,omitempty"`
}

type modelLite struct {
	Properties map[string]*propertyLite `yaml:"properties,omitempty"`
}

func (ff *FileFetcher) selectFieldFromModelFile(input input.IInputMethods, file string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	var m modelLite
	if err := yaml.Unmarshal(b, &m); err != nil {
		return "", fmt.Errorf("[ERROR] failed to parse YAML: %s", file)
	}

	if len(m.Properties) == 0 {
		return "", fmt.Errorf("[ERROR] no properties found in model: %s", file)
	}

	// interactive walk over properties
	pointer := PROPERTIES_PATH
	current := m.Properties

	for {
		keys := ff.sortedKeys(current)
		var sel string
		if err := input.SelectInput(&sel, SELECT_PROPERTY_MSG, keys); err != nil {
			return "", err
		}
		pointer = pointer + "/" + ff.baseFetcher.EscapeJsonPointerToken(sel)
		prop := current[sel]

		// Decide next
		// If object with sub-properties
		if prop != nil && prop.Type == constants.OBJECT_TYPE && len(prop.Properties) > 0 {
			// Ask continue or use here
			var goDeeper string
			if err := input.SelectInput(&goDeeper, CONTINUE_INTO_SUB_PROPS, []string{YES_OPTION, USE_THIS_FIELD}); err != nil {
				return "", err
			}
			if goDeeper == YES_OPTION {
				current = prop.Properties
				pointer = pointer + PROPERTIES_PATH
				continue
			}
			return pointer, nil
		}
		// If array with items
		if prop != nil && prop.Type == constants.ARRAY_TYPE && prop.Items != nil {
			var goItems string
			if err := input.SelectInput(&goItems, ARRAY_ITEMS_OR_USE, []string{ITEMS_OPTION, USE_THIS_FIELD}); err != nil {
				return "", err
			}
			if goItems == ITEMS_OPTION {
				pointer = pointer + ITEMS_PATH
				// dive into the item shape
				if prop.Items.Type == constants.OBJECT_TYPE && len(prop.Items.Properties) > 0 {
					current = prop.Items.Properties
					pointer = pointer + PROPERTIES_PATH
					continue
				}
				// items is primitive or non-object
				return pointer, nil
			}
			return pointer, nil
		}
		// primitive or no deeper structure
		return pointer, nil
	}
}

// selectFieldFromSchemaFile parses a schema file (map root) and guides the user to select a field.
// Returns a JSON Pointer like "/<SchemaName>/properties/foo" (without leading '#').
func (ff *FileFetcher) selectFieldFromSchemaFile(input input.IInputMethods, file string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	// schema files are marshaled as: map[schemaName]*propertyLite
	var root map[string]*propertyLite
	if err := yaml.Unmarshal(b, &root); err != nil {
		return "", fmt.Errorf("[ERROR] failed to parse YAML: %s", file)
	}
	if len(root) == 0 {
		return "", fmt.Errorf("[ERROR] schema file has no root entries: %s", file)
	}

	names := ff.sortedKeysPtrMap(root)
	var schemaName string
	if err := input.SelectInput(&schemaName, SELECT_ROOT_SCHEMA_MSG, names); err != nil {
		return "", err
	}

	pointer := "/" + ff.baseFetcher.EscapeJsonPointerToken(schemaName)
	prop := root[schemaName]
	if prop == nil {
		return "", errors.New("selected root not found")
	}

	// If object, dive into properties; otherwise allow use directly
	if prop.Type == constants.OBJECT_TYPE || prop.Type == constants.ARRAY_TYPE {
		// allow using the root as-is
		var use string
		if err := input.SelectInput(&use, USE_ROOT_SCHEMA_MSG, []string{USE_THIS, SELECT_FIELD}); err != nil {
			return "", err
		}
		if use == USE_THIS {
			return pointer, nil
		}
	}

	// Traverse similar to model
	currentProp := prop
	for {
		if currentProp.Type == constants.OBJECT_TYPE && len(currentProp.Properties) > 0 {
			keys := ff.sortedKeys(currentProp.Properties)
			var sel string
			if err := input.SelectInput(&sel, SELECT_PROPERTY_MSG, keys); err != nil {
				return "", err
			}
			pointer = pointer + PROPERTIES_PATH + "/" + ff.baseFetcher.EscapeJsonPointerToken(sel)
			currentProp = currentProp.Properties[sel]
			continue
		}
		if currentProp.Type == constants.ARRAY_TYPE && currentProp.Items != nil {
			var goItems string
			if err := input.SelectInput(&goItems, ARRAY_DETECTED_MSG, []string{ITEMS_OPTION, USE_THIS_FIELD}); err != nil {
				return "", err
			}
			if goItems == ITEMS_OPTION {
				pointer = pointer + ITEMS_PATH
				currentProp = currentProp.Items
				continue
			}
			return pointer, nil
		}
		// primitive
		return pointer, nil
	}
}

func (ff *FileFetcher) sortedKeys(m map[string]*propertyLite) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (ff *FileFetcher) sortedKeysPtrMap(m map[string]*propertyLite) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
