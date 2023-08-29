package goeasyi18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromJson(t *testing.T) {
	// Create a temp directory for testing
	tempDir := filepath.Join(os.TempDir(), "goeasyi18n_tests_loader_json")
	os.RemoveAll(tempDir)
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Errorf("Unexpected error creating temp dir for the tests: %v", err)
	}

	// Helper function to build paths
	buildPath := func(fileName string) string {
		return filepath.Join(tempDir, fileName)
	}

	// Prepare test files
	os.WriteFile(buildPath("test1.json"), []byte(`[{"Key": "hello", "Default": "Hello"}]`), 0755)
	os.WriteFile(buildPath("test2.json"), []byte(`[{"Key": "world", "Default": "World"}]`), 0755)
	os.WriteFile(buildPath("bad.json"), []byte(`bad json`), 0755)

	t.Run("load single file", func(t *testing.T) {
		strings, err := LoadFromJson(buildPath("test1.json"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 1 || strings[0].Key != "hello" {
			t.Errorf("Unexpected result: %v", strings)
		}
		if strings[0].Default != "Hello" {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("load multiple files", func(t *testing.T) {
		strings, err := LoadFromJson(buildPath("test1.json"), buildPath("test2.json"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 2 {
			t.Errorf("Unexpected result: %v", strings)
		}
		if strings[0].Key != "hello" {
			t.Errorf("Unexpected result: %v", strings)
		}
		if strings[1].Key != "world" {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("load with glob pattern", func(t *testing.T) {
		strings, err := LoadFromJson(buildPath("test*.json"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 2 {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("handle bad json", func(t *testing.T) {
		_, err := LoadFromJson(buildPath("bad.json"))
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("handle no match glob", func(t *testing.T) {
		_, err := LoadFromJson(buildPath("nomatch*.json"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("run translation on parsed strings", func(t *testing.T) {
		i18n := NewI18n(Config{})
		english, err := LoadFromJson(buildPath("test1.json"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		i18n.AddLanguage("en", english)

		result := i18n.Translate("en", "hello", TOptions{})
		expected := "Hello"

		if result != expected {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}
