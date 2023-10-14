package goeasyi18n

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadFromYamlBytes loads a list of TranslateString
// from the provided YAML bytes.
func LoadFromYamlBytes(
	yamlBytes []byte,
) (TranslateStrings, error) {
	// This function uses a bridge, it converts YAML to JSON before
	// parsing it as a TranslateString to avoid inconsistencies

	var parsedYaml any
	err := yaml.Unmarshal(yamlBytes, &parsedYaml)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(parsedYaml)
	if err != nil {
		return nil, err
	}

	var translateStrings TranslateStrings
	err = json.Unmarshal(jsonBytes, &translateStrings)
	if err != nil {
		return nil, err
	}

	return translateStrings, nil
}

// LoadFromYamlString loads a list of TranslateString
// from the provided YAML string.
func LoadFromYamlString(
	yamlString string,
) (TranslateStrings, error) {
	return LoadFromYamlBytes([]byte(yamlString))
}

// LoadFromYamlFiles loads a list of TranslateString from
// one or multiple YAML files, allowing glob patterns
// like "path/to/files/*.yaml".
func LoadFromYamlFiles(
	filesOrGlobs ...string,
) (TranslateStrings, error) {
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

			translateString, err := LoadFromYamlBytes(byteValue)
			if err != nil {
				return nil, err
			}

			allTranslateStrings = append(allTranslateStrings, translateString...)
		}
	}

	return allTranslateStrings, nil
}

// LoadFromYamlFS loads a list of TranslateString from
// one or multiple YAML files located within a provided
// filesystem (fs.FS), allowing glob patterns
// like "path/to/files/*.yaml".
func LoadFromYamlFS(
	fileSystem fs.FS,
	filesOrGlobs ...string,
) (TranslateStrings, error) {
	var allTranslateStrings TranslateStrings

	for _, pattern := range filesOrGlobs {
		matches, err := fs.Glob(fileSystem, pattern)
		if err != nil {
			return nil, err
		}

		if len(matches) == 0 {
			continue
		}

		for _, file := range matches {
			byteValue, err := readFileFromFS(fileSystem, file)
			if err != nil {
				return nil, err
			}

			translateString, err := LoadFromYamlBytes(byteValue)
			if err != nil {
				return nil, err
			}

			allTranslateStrings = append(allTranslateStrings, translateString...)
		}
	}

	return allTranslateStrings, nil
}
