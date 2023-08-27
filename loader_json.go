package goeasyi18n

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Load a list of TranslateString from one or multiple JSON files
// It allow glob patterns like "path/to/files/*.json"
func LoadFromJson(filesOrGlobs ...string) (TranslateStrings, error) {
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

			var translateString TranslateStrings
			err = json.Unmarshal(byteValue, &translateString)
			if err != nil {
				return nil, err
			}

			allTranslateStrings = append(allTranslateStrings, translateString...)
		}
	}

	return allTranslateStrings, nil
}
