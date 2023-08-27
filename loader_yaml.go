package goeasyi18n

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load a list of TranslateString from one or multiple YAML files
// It allow glob patterns like "path/to/files/*.yaml"
func LoadFromYaml(filesOrGlobs ...string) (TranslateStrings, error) {
	var allTranslateStrings TranslateStrings

	for _, pattern := range filesOrGlobs {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, file := range matches {
			byteValue, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}

			// This function uses a bridge, it converts YAML to JSON before
			// parsing it as a TranslateString to avoid inconsistencies

			var parsedYaml any
			err = yaml.Unmarshal(byteValue, &parsedYaml)
			if err != nil {
				return nil, err
			}

			jsonBytes, err := json.Marshal(parsedYaml)
			if err != nil {
				return nil, err
			}

			var translateString TranslateStrings
			err = json.Unmarshal(jsonBytes, &translateString)
			if err != nil {
				return nil, err
			}

			allTranslateStrings = append(allTranslateStrings, translateString...)
		}
	}

	return allTranslateStrings, nil
}
