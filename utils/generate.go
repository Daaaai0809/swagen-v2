package utils

import (
	"os"
	"strings"
)

func GenerateSchema(input []byte, fileName, path string) error {
	var name string
	if strings.HasSuffix(path, "/") {
		name = path + fileName + ".yaml"
	} else {
		name = path + "/" + fileName + ".yaml"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	if err := os.WriteFile(name, input, 0644); err != nil {
		return err
	}

	return nil
}
