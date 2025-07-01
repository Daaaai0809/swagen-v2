package utils

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func GenerateSchema(input interface{}, fileName, path string) error {
	data, err := yaml.Marshal(input)
	if err != nil {
		return err
	}

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

	if err := os.WriteFile(name, data, 0644); err != nil {
		return err
	}

	return nil
}
