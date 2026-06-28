package main

import (
	"log"
	"time"

	"github.com/leonard-atorough/tsdb/internal"
	"github.com/leonard-atorough/tsdb/internal/collector"
	"github.com/leonard-atorough/tsdb/internal/models"
)

func main() {
	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}

func start(directory string, tenantID string, pollingInterval time.Duration) {
	config := &internal.Config{
		DataDir:         directory,
		TenantID:        tenantID,
		PollingInterval: pollingInterval,
	}

	filePath := config.GetFilePath("data")

	wr, err := internal.NewFileWriter(filePath)
	if err != nil {
		log.Fatalf("Error creating file writer: %v", err)
	}
	defer wr.Close()

	collector := collector.NewSystemCollector(pollingInterval, "localhost")

	ticker := time.NewTicker(pollingInterval) // Collect data based on the polling interval
	defer ticker.Stop()

	for range ticker.C {
		data, err := collector.Collect()
		if err != nil {
			log.Printf("Error collecting data: %v", err)
			continue
		}
		for _, d := range data {
			err := wr.WriteData(d)
			if err != nil {
				log.Printf("Error writing data: %v", err)
			}
		}
	}
}

func query(startTime, endTime time.Time) ([]models.TimeSeriesData, error) {
	// Placeholder for query implementation
	// This function should read the data from the file and filter it based on the provided time range.
	return nil, nil
}
