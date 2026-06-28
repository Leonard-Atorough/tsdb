package reader

import (
	"path/filepath"
	"testing"
)

func NewReaderTest(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	_, err := NewReader(testPath)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}
}
