package collector

import (
	"time"

	"github.com/leonard-atorough/tsdb/internal/models"
)

type Collector interface {
	Collect() ([]models.TimeSeriesData, error)
}

type MockCollector struct{}

func (mc *MockCollector) Collect() ([]models.TimeSeriesData, error) {
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
	metrics  *SystemMetrics
	hostname string
}

func NewSystemCollector(cpuDuration time.Duration, hostname string) *SystemCollector {
	return &SystemCollector{
		metrics:  NewSystemMetrics(cpuDuration),
		hostname: hostname,
	}
}

func (sc *SystemCollector) Collect() ([]models.TimeSeriesData, error) {
	var data []models.TimeSeriesData

	cpuUsage, err := sc.metrics.CPUUsage()
	if err != nil {
		return nil, err
	}

	if sc.hostname == "" {
		sc.hostname = "localhost"
	}
	
	data = append(data, models.TimeSeriesData{
		Measurement: "cpu",
		TagSet: map[string]string{
			"host": sc.hostname,
		},
		FieldSet: map[string]any{
			"usage": cpuUsage,
		},
		Timestamp: time.Now().UnixMilli(),
	})

	memUsage, err := sc.metrics.MemoryUsage()
	if err != nil {
		return nil, err
	}
	data = append(data, models.TimeSeriesData{
		Measurement: "memory",
		TagSet: map[string]string{
			"host": sc.hostname,
		},
		FieldSet: map[string]any{
			"usage": memUsage,
		},
		Timestamp: time.Now().UnixMilli(),
	})

	return data, nil
}
