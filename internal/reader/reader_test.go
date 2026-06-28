package reader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leonard-atorough/tsdb/internal"
	"github.com/stretchr/testify/assert"
)

func generateTestData(t *testing.T, filepath string, data []string) {
	content := strings.Join(data, "\n")
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

}

func testData() []string {
	return []string{
		`{"m":"cpu","f":{"value":50.0},"t":{"host":"a","env":"prod"},"ts":1000}`,
		`{"m":"cpu","f":{"value":60.0},"t":{"host":"a","env":"staging"},"ts":2000}`,
		`{"m":"cpu","f":{"value":70.0},"t":{"host":"b","env":"prod"},"ts":3000}`,
		`invalid json line`, // Tests graceful skip
		`{}`,                // Edge case: empty object
	}
}

func NewReaderTest(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	_, err := NewReader(testPath)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}
}

func TestReader_Query_ValidTimeRange(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, err := NewReader(testPath)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000)})

	assert.Len(t, points, 4, "Expected 4 valid points in the time range")
	assert.Equal(t, "cpu", points[0].Measurement)
}

func TestReader_Query_InvalidTimeRange(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, err := NewReader(testPath)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(5000), To: internal.ConvertUnixToTime(6000)})

	assert.Len(t, points, 0, "Expected 0 points in the invalid time range")
}
