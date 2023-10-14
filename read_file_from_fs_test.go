package goeasyi18n

import (
	"embed"
	"testing"
)

//go:embed testfiles/*
var readFSTestFiles embed.FS

func TestReadFileFromFS(t *testing.T) {
	tests := []struct {
		filename     string
		expectedData string
		expectedErr  bool
	}{
		{
			filename:     "testfiles/test1.txt",
			expectedData: "test 1",
			expectedErr:  false,
		},
		{
			filename:     "testfiles/test2.txt",
			expectedData: "test 2",
			expectedErr:  false,
		},
		{
			filename:     "testfiles/not-exists.txt",
			expectedData: "",
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			data, err := readFileFromFS(readFSTestFiles, tt.filename)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("didn't expect error, got %v", err)
				}
			}

			if string(data) != tt.expectedData {
				t.Errorf("expected data %q, got %q", tt.expectedData, string(data))
			}
		})
	}
}
