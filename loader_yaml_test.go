package goeasyi18n

import (
	"embed"
	"testing"
)

func TestLoadFromYamlBytes(t *testing.T) {
	bytes := []byte("- Key: hello\n  Default: Hello\n")
	strings, err := LoadFromYamlBytes(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(strings) != 1 || strings[0].Key != "hello" {
		t.Errorf("Unexpected result: %v", strings)
	}
	if strings[0].Default != "Hello" {
		t.Errorf("Unexpected result: %v", strings)
	}
}

func TestLoadFromYamlString(t *testing.T) {
	strings, err := LoadFromYamlString("- Key: hello\n  Default: Hello\n")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(strings) != 1 || strings[0].Key != "hello" {
		t.Errorf("Unexpected result: %v", strings)
	}
	if strings[0].Default != "Hello" {
		t.Errorf("Unexpected result: %v", strings)
	}
}

func TestLoadFromYamlFiles(t *testing.T) {
	t.Run("load single file", func(t *testing.T) {
		strings, err := LoadFromYamlFiles("./testfiles/test1.yaml")
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
		strings, err := LoadFromYamlFiles(
			"./testfiles/test1.yaml",
			"./testfiles/test2.yaml",
		)
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
		strings, err := LoadFromYamlFiles("./testfiles/test*.yaml")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 2 {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("handle incorrect yaml", func(t *testing.T) {
		_, err := LoadFromYamlFiles("./testfiles/incorrect.yaml")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("handle no match glob", func(t *testing.T) {
		_, err := LoadFromYamlFiles("./testfiles/nomatch*.yaml")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

//go:embed testfiles/*
var yamlTestFiles embed.FS

func TestLoadFromYamlFS(t *testing.T) {
	t.Run("load single file", func(t *testing.T) {
		strings, err := LoadFromYamlFS(
			yamlTestFiles,
			"testfiles/test1.yaml",
		)
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
		strings, err := LoadFromYamlFS(
			yamlTestFiles,
			"testfiles/test1.yaml",
			"testfiles/test2.yaml",
		)
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
		strings, err := LoadFromYamlFS(
			yamlTestFiles,
			"testfiles/test*.yaml",
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 2 {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("handle incorrect yaml", func(t *testing.T) {
		_, err := LoadFromYamlFS(
			yamlTestFiles,
			"testfiles/incorrect.yaml",
		)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("handle no match glob", func(t *testing.T) {
		_, err := LoadFromYamlFS(
			yamlTestFiles,
			"testfiles/nomatch*.yaml",
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
