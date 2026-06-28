package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leonard-atorough/tsdb/internal/models"
)

func TestNewFileWriter(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	fw, err := NewFileWriter(testPath)

	if err != nil {
		t.Fatalf("Failed to create FileWriter: %v", err)
	}
	if fw == nil {
		t.Errorf("NewFileWriter returned nil")
	}
	if fw.path == "" {
		t.Errorf("FileWriter path is empty")
	}

	fw.Close()
}

func TestWriteData(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_write.jsonl")

	fw, err := NewFileWriter(testPath)
	if err != nil {
		t.Fatalf("Failed to create FileWriter: %v", err)
	}
	defer fw.Close()

	data := models.TimeSeriesData{
		Measurement: "cpu",
		TagSet: map[string]string{
			"host": "server1",
		},
		FieldSet: map[string]any{
			"usage": 75.5,
		},
		Timestamp: 1234567890,
	}

	err = fw.WriteData(data)

	if err != nil {
		t.Errorf("WriteData failed: %v", err)
	}

	fw.Close()

	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	if len(content) == 0 {
		t.Errorf("File is empty after WriteData")
	}
}

func TestWriteDataWithClosedFile(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_write_fail.jsonl")

	fw, err := NewFileWriter(testPath)
	if err != nil {
		t.Fatalf("Failed to create FileWriter: %v", err)
	}
	fw.Close()

	data := models.TimeSeriesData{
		Measurement: "cpu",
		FieldSet: map[string]any{
			"usage": 75.5,
		},
		Timestamp: 1234567890,
	}

	err = fw.WriteData(data)

	if err == nil {
		t.Errorf("Expected error writing to closed file, but succeeded")
	}
}

func TestClose(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_close.jsonl")

	testFile, err := os.Create(testPath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fw := &FileWriter{
		file: testFile,
		path: testPath,
	}

	err = fw.Close()

	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	_, err = fw.file.WriteString("test")
	if err == nil {
		t.Errorf("Expected error writing to closed file, but got none")
	}
}

func TestWriteMultipleRecords(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_multiple.jsonl")

	fw, err := NewFileWriter(testPath)
	if err != nil {
		t.Fatalf("Failed to create FileWriter: %v", err)
	}
	defer fw.Close()

	records := []models.TimeSeriesData{
		{
			Measurement: "cpu",
			FieldSet: map[string]any{
				"usage": 50.0,
			},
			Timestamp: 1000,
		},
		{
			Measurement: "cpu",
			FieldSet: map[string]any{
				"usage": 60.0,
			},
			Timestamp: 1001,
		},
		{
			Measurement: "cpu",
			FieldSet: map[string]any{
				"usage": 70.0,
			},
			Timestamp: 1002,
		},
	}

	for _, data := range records {
		err := fw.WriteData(data)
		if err != nil {
			t.Errorf("WriteData failed: %v", err)
		}
	}

	fw.Close()

	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	if len(content) == 0 {
		t.Errorf("File is empty after writing multiple records")
	}
}
