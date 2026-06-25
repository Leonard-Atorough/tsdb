package internal

import (
	"time"

	"github.com/leonard-atorough/tsdb/internal/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Collector interface {
	Collect() ([]models.TimeSeriesData, error)
}

type MockCollector struct {} // MockCollector is a mock implementation of the Collector interface for testing purposes.

func (mc *MockCollector) Collect() ([]models.TimeSeriesData, error) {
	// Return some mock data for testing purposes.
	mockData := []models.TimeSeriesData{
		{
			Measurement: "cpu",
			TagSet: map[string]string{
				"host": "localhost",
			},
			FieldSet: map[string]any{
				"usage": 0.5,
			},
			Timestamp: time.Now().UnixMilli(),
		},
		{
			Measurement: "memory",
			TagSet: map[string]string{
				"host": "localhost",
			},
			FieldSet: map[string]any{
				"usage": 0.75,
			},
			Timestamp: time.Now().UnixMilli(),
		},
	}
	return mockData, nil
}

type SystemCollector struct {
	CPUDuration time.Duration // Duration for CPU usage measurement
}

func NewSystemCollector(cpuDuration time.Duration) *SystemCollector {
	return &SystemCollector{
		CPUDuration: cpuDuration,
	}
}

func (sc *SystemCollector) Collect() ([]models.TimeSeriesData, error) {
	processMetrics := []models.TimeSeriesData{}

	cpuPercentages, err := cpu.Percent(sc.CPUDuration, false)
	if err != nil {
		return nil, err
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	cpuData := models.TimeSeriesData{
		Measurement: "cpu",
		TagSet: map[string]string{
			"host": "localhost",
		},
		FieldSet: map[string]any{
			"usage": cpuPercentages[0],
		},
		Timestamp: time.Now().UnixMilli(),
	}
	processMetrics = append(processMetrics, cpuData)

	memoryData := models.TimeSeriesData{
		Measurement: "memory",
		TagSet: map[string]string{
			"host": "localhost",
		},
		FieldSet: map[string]any{
			"usage": vmStat.UsedPercent,
		},
		Timestamp: time.Now().UnixMilli(),
	}
	processMetrics = append(processMetrics, memoryData)

	return processMetrics, nil
}