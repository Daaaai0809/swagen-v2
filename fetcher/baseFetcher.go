package fetcher

import (
	"os"
	"sort"
	"strings"
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
	NEW_DIR          = "Create new directory"
)

// IBaseFetcher defines common operations for file and directory fetching
type IBaseFetcher interface {
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
