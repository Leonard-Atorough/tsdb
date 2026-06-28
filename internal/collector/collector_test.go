package collector

import (
	"testing"
)

func TestMockCollector(t *testing.T) {
	mockCollector := &MockCollector{}

	mockData, err := mockCollector.Collect()
	if err != nil {
		t.Fatalf("MockCollector.Collect() returned an error: %v", err)
	}

	if len(mockData) != 2 {
		t.Fatalf("Expected 2 data points, got %d", len(mockData))
	}

	if mockData[0].Measurement != "cpu" || mockData[1].Measurement != "memory" {
		t.Fatalf("Unexpected measurement names: %s, %s", mockData[0].Measurement, mockData[1].Measurement)
	}
}
