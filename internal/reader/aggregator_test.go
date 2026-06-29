package reader

import (
	"path/filepath"
	"testing"

	"github.com/leonard-atorough/tsdb/internal"
	"github.com/stretchr/testify/assert"
)

func TestAggregator_Average(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	average, _ := r.Avg(AggregateOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Measurement: "cpu", Field: "value"})

	assert.Equal(t, 60.0, average, "Expected average value to be 60.0")
}

func TestAggregator_Sum(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	sum, _ := r.Sum(AggregateOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Measurement: "cpu", Field: "value"})

	assert.Equal(t, 180.0, sum, "Expected sum value to be 180.0")
}

func TestAggregator_Count(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	count, _ := r.Count(AggregateOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Measurement: "cpu", Field: "value"})

	assert.Equal(t, 3, count, "Expected count value to be 3")
}

func TestAggregator_Min(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	min, _ := r.Min(AggregateOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Measurement: "cpu", Field: "value"})

	assert.Equal(t, 50.0, min, "Expected min value to be 50.0")
}

func TestAggregator_Max(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	max, _ := r.Max(AggregateOpts{From: internal.ConvertUnixToTime(0), To: internal.ConvertUnixToTime(4000), Measurement: "cpu", Field: "value"})

	assert.Equal(t, 70.0, max, "Expected max value to be 70.0")
}

func TestAggregator_Aggregates(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_file.jsonl")

	generateTestData(t, testPath, testData())
	r, _ := NewReader(testPath)

	aggregates, _ := r.Aggregates(
		AggregateOpts{
			From:        internal.ConvertUnixToTime(0),
			To:          internal.ConvertUnixToTime(4000),
			Measurement: "cpu",
			Field:       "value",
			Funcs:       []string{"sum", "count", "avg", "min", "max"}})

	assert.Equal(t, 180.0, aggregates["sum"], "Expected sum value to be 180.0")
	assert.Equal(t, 3.0, aggregates["count"], "Expected count value to be 3.0")
	assert.Equal(t, 60.0, aggregates["avg"], "Expected average value to be 60.0")
	assert.Equal(t, 50.0, aggregates["min"], "Expected min value to be 50.0")
	assert.Equal(t, 70.0, aggregates["max"], "Expected max value to be 70.0")
}
