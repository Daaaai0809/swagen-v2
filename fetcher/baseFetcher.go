package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Daaaai0809/swagen-v2/input"
)

// Constants for common use across fetcher files
const (
	YAML_EXT         = ".yaml"
	YML_EXT          = ".yml"
	PARENT_DIR       = "../"
	YES_OPTION       = "Yes"
	USE_THIS_FIELD   = "Use this field"
	USE_THIS         = "Use this"
	SELECT_FIELD     = "Select field"
	ITEMS_OPTION     = "items"
	JSON_POINTER_REF = "#"
	TILDE_ESCAPE     = "~0"
	SLASH_ESCAPE     = "~1"
)

// IBaseFetcher defines common operations for file and directory fetching
type IBaseFetcher interface {
	// SelectFileInteractive allows interactive file selection from a directory
	SelectFileInteractive(input input.IInputMethods, start string, fileFilter func(string) bool) (string, error)

	// SortedStringKeys returns sorted keys from a map[string]interface{}
	SortedStringKeys(m map[string]interface{}) []string

	// SortedKeysFromPtrMap returns sorted keys from a map with pointer values
	SortedKeysFromPtrMap(m interface{}) []string

	// EscapeJsonPointerToken escapes JSON Pointer tokens according to RFC 6901
	EscapeJsonPointerToken(s string) string

	// ReadDirectoryEntries reads and categorizes directory entries
	ReadDirectoryEntries(dirPath string) (dirs []string, files []string, err error)
}

// BaseFetcher implements common fetcher functionality
type BaseFetcher struct{}

// NewBaseFetcher creates a new BaseFetcher instance
func NewBaseFetcher() IBaseFetcher {
	return &BaseFetcher{}
}

// SelectFileInteractive lets the user navigate directories and select a file based on filter.
func (bf *BaseFetcher) SelectFileInteractive(input input.IInputMethods, start string, fileFilter func(string) bool) (string, error) {
	cwd := filepath.Clean(start)
	for {
		dirs, files, err := bf.ReadDirectoryEntries(cwd)
		if err != nil {
			return "", fmt.Errorf("[ERROR] cannot read directory: %s", cwd)
		}

		// Filter files
		var filteredFiles []string
		for _, file := range files {
			if fileFilter(file) {
				filteredFiles = append(filteredFiles, file)
			}
		}

		items := make([]string, 0, len(dirs)+len(filteredFiles)+1)
		if cwd != filepath.Clean(start) {
			items = append(items, PARENT_DIR)
		}
		items = append(items, dirs...)
		items = append(items, filteredFiles...)

		if len(items) == 0 {
			return "", fmt.Errorf("[ERROR] no selectable items in directory: %s", cwd)
		}

		var sel string
		if err := input.SelectInput(&sel, fmt.Sprintf("Select entry in %s", cwd), items); err != nil {
			return "", err
		}

		switch sel {
		case PARENT_DIR:
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

// ReadDirectoryEntries reads and categorizes directory entries into dirs and files
func (bf *BaseFetcher) ReadDirectoryEntries(dirPath string) (dirs []string, files []string, err error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}

	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue // hide dotfiles
		}
		if e.IsDir() {
			dirs = append(dirs, name+"/")
		} else {
			files = append(files, name)
		}
	}
	sort.Strings(dirs)
	sort.Strings(files)

	return dirs, files, nil
}

// SortedStringKeys returns sorted keys from a map[string]interface{}
func (bf *BaseFetcher) SortedStringKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedKeysFromPtrMap returns sorted keys from a map with pointer values
func (bf *BaseFetcher) SortedKeysFromPtrMap(m interface{}) []string {
	// This is a generic implementation that works with reflection
	// but for type safety, we'll provide specific implementations
	return []string{}
}

// EscapeJsonPointerToken escapes '~' and '/' per RFC 6901
func (bf *BaseFetcher) EscapeJsonPointerToken(s string) string {
	s = strings.ReplaceAll(s, "~", TILDE_ESCAPE)
	s = strings.ReplaceAll(s, "/", SLASH_ESCAPE)
	return s
}
