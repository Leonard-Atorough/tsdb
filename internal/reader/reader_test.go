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
		`{"m":"cpu","fields":{"value":50.0},"tags":{"host":"a","env":"prod"},"ts":1000}`,
		`{"m":"cpu","fields":{"value":60.0},"tags":{"host":"a","env":"staging"},"ts":2000}`,
		`{"m":"cpu","fields":{"value":70.0},"tags":{"host":"b","env":"prod"},"ts":3000}`,
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

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000)})

	assert.Len(t, points, 4, "Expected 4 valid points in the time range")
	assert.Equal(t, "cpu", points[0].Measurement)
}

func TestReader_Query_ValidTimeRangeWithNoEndTime(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0)})

	assert.Len(t, points, 4, "Expected 4 valid points in the time range")
	assert.Equal(t, "cpu", points[0].Measurement)
}

func TestReader_Query_ValidTimeRangeWithNoStartTime(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{To: internal.ConvertUnixToTime(4000)})

	assert.Len(t, points, 4, "Expected 4 valid points in the time range")
	assert.Equal(t, "cpu", points[0].Measurement)
}

func TestReader_Query_ValidTimeRangeWithAgo(t *testing.T) {
	// testing the "Ago" functionality by setting a time range of 3 seconds ago
	// validate that from is set to 3 seconds before the current time and to is set to the current time
	
	r, _ := setupTestReader(t, testData())

	// Don't bother testing data, just that the time range is set correctly

}

func TestReader_Query_InvalidTimeRange(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(5000), To: internal.ConvertUnixToTime(6000)})

	assert.Len(t, points, 0, "Expected 0 points in the invalid time range")
}

func TestReader_Query_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "empty_file.jsonl")

	// Create an empty file
	if err := os.WriteFile(testPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(1000)})

	assert.Len(t, points, 0, "Expected 0 points in the empty file")
}

func TestReader_Query_WithMeasurementFilter(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Measurement: "cpu"})

	assert.Len(t, points, 3, "Expected 3 points with measurement 'cpu'")
}

func TestReader_Query_WithTagFilter(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Tags: map[string]string{"host": "a"}})

	assert.Len(t, points, 2, "Expected 2 points with tag host='a'")
}

func TestReader_Query_WithMultipleTags(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())

	r, _ := NewReader(testPath)

	points, _ := r.Query(QueryOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Tags: map[string]string{"host": "a", "env": "staging"}})

	assert.Len(t, points, 1, "Expected 1 point with tags host='a' and env='staging'")
}

func setupTestReader(t *testing.T, data []string) (*Reader, string) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, data)
	r, _ := NewReader(testPath)
	return r, testPath
}