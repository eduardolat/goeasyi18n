package goeasyi18n

import (
	"embed"
	"testing"
)

func TestLoadFromJsonBytes(t *testing.T) {
	bytes := []byte(`[{"Key": "hello", "Default": "Hello"}]`)
	strings, err := LoadFromJsonBytes(bytes)
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

func TestLoadFromJsonString(t *testing.T) {
	strings, err := LoadFromJsonString(`[{"Key": "hello", "Default": "Hello"}]`)
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

func TestLoadFromJsonFiles(t *testing.T) {
	t.Run("load single file", func(t *testing.T) {
		strings, err := LoadFromJsonFiles("./testfiles/test1.json")
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
		strings, err := LoadFromJsonFiles(
			"./testfiles/test1.json",
			"./testfiles/test2.json",
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
		strings, err := LoadFromJsonFiles("./testfiles/test*.json")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 2 {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("handle incorrect json", func(t *testing.T) {
		_, err := LoadFromJsonFiles("./testfiles/incorrect.json")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("handle no match glob", func(t *testing.T) {
		_, err := LoadFromJsonFiles("./testfiles/nomatch*.json")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

//go:embed testfiles/*
var jsonTestFiles embed.FS

func TestLoadFromJsonFS(t *testing.T) {
	t.Run("load single file", func(t *testing.T) {
		strings, err := LoadFromJsonFS(
			jsonTestFiles,
			"testfiles/test1.json",
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
		strings, err := LoadFromJsonFS(
			jsonTestFiles,
			"testfiles/test1.json",
			"testfiles/test2.json",
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
		strings, err := LoadFromJsonFS(
			jsonTestFiles,
			"testfiles/test*.json",
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(strings) != 2 {
			t.Errorf("Unexpected result: %v", strings)
		}
	})

	t.Run("handle incorrect json", func(t *testing.T) {
		_, err := LoadFromJsonFS(
			jsonTestFiles,
			"testfiles/incorrect.json",
		)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("handle no match glob", func(t *testing.T) {
		_, err := LoadFromJsonFS(
			jsonTestFiles,
			"testfiles/nomatch*.json",
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
