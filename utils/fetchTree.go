package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Daaaai0809/swagen-v2/constants"
	"gopkg.in/yaml.v2"
)

// InteractiveResolveRef
// Contract:
// - Inputs: IInputMethods, mode (SCHEMA or API)
// - Outcome: returns a $ref string built from a user-selected YAML file and field
// - Behavior:
//   - SCHEMA mode: traverse from MODEL_PATH
//   - API mode: ask user to choose start path (MODEL_PATH or SCHEMA_PATH)
//   - After selecting a file, parse YAML into Model or Schema(map) and select a field
//   - Build JSON Pointer and relative file path (relative to SCHEMA_PATH or API_PATH)
func InteractiveResolveRef(input IInputMethods, mode constants.InputMode) (string, error) {
	// Decide start path and destination base to compute relative path
	startPath, fileKind, err := decideStartPath(input, mode)
	if err != nil {
		return "", err
	}

	destBase := GetEnv(SCHEMA_PATH, "")
	if mode == constants.MODE_API {
		destBase = GetEnv(API_PATH, "")
	}
	if destBase == "" {
		destBase = "."
	}

	// Directory traversal to pick a YAML file
	selectedFile, err := selectYamlFileInteractive(input, startPath)
	if err != nil {
		return "", err
	}

	// If startPath came from SCHEMA_PATH, treat file as schema-kind even in API mode
	if fileKind == "auto" {
		// infer by extension only (already .yaml) and location
		if strings.HasPrefix(filepath.Clean(selectedFile), filepath.Clean(GetEnv(SCHEMA_PATH, ""))) {
			fileKind = "schema"
		} else {
			fileKind = "model"
		}
	}

	// Parse YAML and select a field to build JSON Pointer
	var pointer string
	switch fileKind {
	case "model":
		pointer, err = selectFieldFromModelFile(input, selectedFile)
	case "schema":
		pointer, err = selectFieldFromSchemaFile(input, selectedFile)
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
	if !strings.HasSuffix(rel, ".yaml") && !strings.HasSuffix(rel, ".yml") {
		// safety: ensure extension
		rel += ".yaml"
	}

	// Ensure pointer starts with '#'
	if pointer == "" {
		pointer = "#"
	} else if !strings.HasPrefix(pointer, "#") {
		pointer = "#" + pointer
	}

	return fmt.Sprintf("%s%s", rel, pointer), nil
}

// decideStartPath asks for start directory based on mode and returns also the fileKind hint
// fileKind: "model", "schema", or "auto"
func decideStartPath(input IInputMethods, mode constants.InputMode) (string, string, error) {
	switch mode {
	case constants.MODE_SCHEMA:
		p := GetEnv(MODEL_PATH, "")
		if p == "" {
			return "", "", errors.New("[ERROR] MODEL_PATH is not set. Set it in environment or .env")
		}
		return p, "model", nil
	case constants.MODE_API:
		var choice string
		if err := input.SelectInput(&choice, "[INFO] Which base path to reference?", []string{"MODEL_PATH", "SCHEMA_PATH"}); err != nil {
			return "", "", err
		}
		switch choice {
		case "MODEL_PATH":
			p := GetEnv(MODEL_PATH, "")
			if p == "" {
				return "", "", errors.New("[ERROR] MODEL_PATH is not set. Set it in environment or .env")
			}
			return p, "model", nil
		case "SCHEMA_PATH":
			p := GetEnv(SCHEMA_PATH, "")
			if p == "" {
				return "", "", errors.New("[ERROR] SCHEMA_PATH is not set. Set it in environment or .env")
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
func selectYamlFileInteractive(input IInputMethods, start string) (string, error) {
	cwd := filepath.Clean(start)
	for {
		entries, err := os.ReadDir(cwd)
		if err != nil {
			return "", fmt.Errorf("[ERROR] cannot read directory: %s", cwd)
		}

		var dirs, files []string
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, ".") {
				continue // hide dotfiles
			}
			if e.IsDir() {
				dirs = append(dirs, name+"/")
			} else if strings.HasSuffix(strings.ToLower(name), ".yaml") || strings.HasSuffix(strings.ToLower(name), ".yml") {
				files = append(files, name)
			}
		}
		sort.Strings(dirs)
		sort.Strings(files)

		items := make([]string, 0, len(dirs)+len(files)+1)
		if cwd != filepath.Clean(start) {
			items = append(items, "../")
		}
		items = append(items, dirs...)
		items = append(items, files...)

		if len(items) == 0 {
			return "", fmt.Errorf("[ERROR] no selectable items in directory: %s", cwd)
		}

		var sel string
		if err := input.SelectInput(&sel, fmt.Sprintf("[INFO] Select entry in %s", cwd), items); err != nil {
			return "", err
		}

		switch sel {
		case "../":
			cwd = filepath.Dir(cwd)
			continue
		default:
			// directory?
			if strings.HasSuffix(sel, "/") {
				cwd = filepath.Join(cwd, strings.TrimSuffix(sel, "/"))
				continue
			}
			// file selected
			return filepath.Join(cwd, sel), nil
		}
	}
}

// selectFieldFromModelFile parses a model file and guides the user to select a field.
// Returns a JSON Pointer like "/properties/foo/items/properties/bar" (without leading '#').
// local lite types to avoid importing libs and causing cycles
type propertyLite struct {
	Type       string                   `yaml:"type,omitempty"`
	Properties map[string]*propertyLite `yaml:"properties,omitempty"`
	Items      *propertyLite            `yaml:"items,omitempty"`
}

type modelLite struct {
	Properties map[string]*propertyLite `yaml:"properties,omitempty"`
}

func selectFieldFromModelFile(input IInputMethods, file string) (string, error) {
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
	pointer := "/properties"
	current := m.Properties

	for {
		keys := sortedKeys(current)
		var sel string
		if err := input.SelectInput(&sel, "[INFO] Select property", keys); err != nil {
			return "", err
		}
		pointer = pointer + "/" + escapeJsonPointerToken(sel)
		prop := current[sel]

		// Decide next
		// If object with sub-properties
		if prop != nil && prop.Type == constants.OBJECT_TYPE && len(prop.Properties) > 0 {
			// Ask continue or use here
			var goDeeper string
			if err := input.SelectInput(&goDeeper, "[INFO] This is an object. Continue into sub-properties?", []string{"Yes", "Use this field"}); err != nil {
				return "", err
			}
			if goDeeper == "Yes" {
				current = prop.Properties
				pointer = pointer + "/properties"
				continue
			}
			return pointer, nil
		}
		// If array with items
		if prop != nil && prop.Type == constants.ARRAY_TYPE && prop.Items != nil {
			var goItems string
			if err := input.SelectInput(&goItems, "[INFO] This is an array. Select items or use array as is?", []string{"items", "Use this field"}); err != nil {
				return "", err
			}
			if goItems == "items" {
				pointer = pointer + "/items"
				// dive into the item shape
				if prop.Items.Type == constants.OBJECT_TYPE && len(prop.Items.Properties) > 0 {
					current = prop.Items.Properties
					pointer = pointer + "/properties"
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
func selectFieldFromSchemaFile(input IInputMethods, file string) (string, error) {
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

	names := sortedKeysPtrMap(root)
	var schemaName string
	if err := input.SelectInput(&schemaName, "[INFO] Select root schema", names); err != nil {
		return "", err
	}

	pointer := "/" + escapeJsonPointerToken(schemaName)
	prop := root[schemaName]
	if prop == nil {
		return "", errors.New("selected root not found")
	}

	// If object, dive into properties; otherwise allow use directly
	if prop.Type == constants.OBJECT_TYPE || prop.Type == constants.ARRAY_TYPE {
		// allow using the root as-is
		var use string
		if err := input.SelectInput(&use, "[INFO] Use the root schema as reference?", []string{"Use this", "Select field"}); err != nil {
			return "", err
		}
		if use == "Use this" {
			return pointer, nil
		}
	}

	// Traverse similar to model
	currentProp := prop
	for {
		if currentProp.Type == constants.OBJECT_TYPE && len(currentProp.Properties) > 0 {
			keys := sortedKeys(currentProp.Properties)
			var sel string
			if err := input.SelectInput(&sel, "[INFO] Select property", keys); err != nil {
				return "", err
			}
			pointer = pointer + "/properties/" + escapeJsonPointerToken(sel)
			currentProp = currentProp.Properties[sel]
			continue
		}
		if currentProp.Type == constants.ARRAY_TYPE && currentProp.Items != nil {
			var goItems string
			if err := input.SelectInput(&goItems, "[INFO] Array detected. Select items or use array as is?", []string{"items", "Use this field"}); err != nil {
				return "", err
			}
			if goItems == "items" {
				pointer = pointer + "/items"
				currentProp = currentProp.Items
				continue
			}
			return pointer, nil
		}
		// primitive
		return pointer, nil
	}
}

func sortedKeys(m map[string]*propertyLite) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysPtrMap[T any](m map[string]*T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// escapeJsonPointerToken escapes '~' and '/' per RFC 6901
func escapeJsonPointerToken(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}
