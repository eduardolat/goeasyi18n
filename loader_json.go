package goeasyi18n

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

// LoadFromJsonBytes loads a list of TranslateString
// from the provided JSON bytes.
func LoadFromJsonBytes(
	jsonBytes []byte,
) (TranslateStrings, error) {
	var translateStrings TranslateStrings
	err := json.Unmarshal(jsonBytes, &translateStrings)
	if err != nil {
		return nil, err
	}

	return translateStrings, nil
}

// LoadFromJsonString loads a list of TranslateString
// from the provided JSON string.
func LoadFromJsonString(
	jsonString string,
) (TranslateStrings, error) {
	return LoadFromJsonBytes([]byte(jsonString))
}

// LoadFromJsonFiles loads a list of TranslateString from
// one or multiple JSON files, allowing glob patterns
// like "path/to/files/*.json".
func LoadFromJsonFiles(
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

			translateString, err := LoadFromJsonBytes(byteValue)
			if err != nil {
				return nil, err
			}

			allTranslateStrings = append(allTranslateStrings, translateString...)
		}
	}

	return allTranslateStrings, nil
}

// LoadFromJsonFS loads a list of TranslateString from
// one or multiple JSON files located within a provided
// filesystem (fs.FS), allowing glob patterns
// like "path/to/files/*.json".
func LoadFromJsonFS(
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

			translateString, err := LoadFromJsonBytes(byteValue)
			if err != nil {
				return nil, err
			}

			allTranslateStrings = append(allTranslateStrings, translateString...)
		}
	}

	return allTranslateStrings, nil
}
