package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Daaaai0809/swagen-v2/constants"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/Daaaai0809/swagen-v2/validator"
)

const (
	INPUT_LABEL_CREATE_DIR = "Enter new directory name: "
)

type IDirectoryFetcher interface {
	InteractiveResolveDir(input input.IInputMethods, mode constants.InputMode) (string, error)
}

type DirectoryFetcher struct {
	BaseFetcher    IBaseFetcher
	InputMethod    input.IInputMethods
	ValidationFunc validator.IInputValidator
}

func NewDirectoryFetcher(inputMethod input.IInputMethods, validationFunc validator.IInputValidator) IDirectoryFetcher {
	return &DirectoryFetcher{
		BaseFetcher:    NewBaseFetcher(),
		InputMethod:    inputMethod,
		ValidationFunc: validationFunc,
	}
}

func (df *DirectoryFetcher) InteractiveResolveDir(input input.IInputMethods, mode constants.InputMode) (string, error) {
	startPath, err := df.decideStartPath(mode)
	if err != nil {
		return "", err
	}

	selectedDir, err := df.selectDirectoryInteractive(input, startPath)
	if err != nil {
		return "", err
	}

	return selectedDir, nil
}

func (df *DirectoryFetcher) selectDirectoryInteractive(input input.IInputMethods, start string) (string, error) {
	cwd := filepath.Clean(start)
	for {
		dirs, _, err := df.BaseFetcher.ReadDirectoryEntries(cwd)
		if err != nil {
			return "", fmt.Errorf("[ERROR] cannot read directory: %s", cwd)
		}

		items := make([]string, 0, len(dirs)+1)
		if cwd != filepath.Clean(start) {
			items = append(items, PARENT_DIR)
		}
		items = append(items, dirs...)
		items = append(items, NEW_DIR)

		if len(items) == 0 {
			return "", fmt.Errorf("[ERROR] no selectable directories in: %s", cwd)
		}

		var sel string
		if err := input.SelectInput(&sel, fmt.Sprintf("Select directory in %s", cwd), items); err != nil {
			return "", err
		}

		switch sel {
		case PARENT_DIR:
			cwd = filepath.Dir(cwd)
			continue
		case NEW_DIR:
			var newDirName string
			if err := df.InputMethod.StringInput(&newDirName, INPUT_LABEL_CREATE_DIR, df.ValidationFunc.Validator_Alphanumeric_Underscore()); err != nil {
				return "", err
			}
			newDirPath := filepath.Join(cwd, newDirName)
			if err := df.createDirectoryIfNotExists(newDirPath); err != nil {
				return "", err
			}
			return newDirPath, nil
		default:
			// directory?
			if strings.HasSuffix(sel, "/") {
				cwd = filepath.Join(cwd, strings.TrimSuffix(sel, "/"))
				// Confirm selection
				confirmItems := []string{USE_THIS, "Continue browsing"}
				var confirmSel string
				if err := input.SelectInput(&confirmSel, fmt.Sprintf("Selected directory: %s. What next?", cwd), confirmItems); err != nil {
					return "", err
				}
				if confirmSel == USE_THIS {
					return cwd, nil
				}
				continue
			}
			// Should not reach here as only directories are listed
			return "", fmt.Errorf("[ERROR] invalid selection: %s", sel)
		}
	}
}

func (df *DirectoryFetcher) decideStartPath(mode constants.InputMode) (string, error) {
	switch mode {
	case constants.MODE_MODEL:
		p := utils.GetEnv(utils.SWAGEN_MODEL_PATH, "")
		if p == "" {
			return "", fmt.Errorf("[ERROR] model path not set")
		}
		return p, nil
	case constants.MODE_SCHEMA:
		p := utils.GetEnv(utils.SWAGEN_SCHEMA_PATH, "")
		if p == "" {
			return "", fmt.Errorf("[ERROR] schema path not set")
		}
		return p, nil
	case constants.MODE_API:
		p := utils.GetEnv(utils.SWAGEN_API_PATH, "")
		if p == "" {
			return "", fmt.Errorf("[ERROR] api path not set")
		}
		return p, nil
	default:
		return "", fmt.Errorf("[ERROR] unsupported mode: %s", mode)
	}
}

func (df *DirectoryFetcher) createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("[ERROR] failed to create directory: %s", path)
		}

		fmt.Printf("[INFO] Created directory: %s\n", path)

		return nil
	}

	fmt.Printf("[INFO] Directory %s is already exists", path)

	return nil
}
