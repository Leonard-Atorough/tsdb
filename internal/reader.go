package internal

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/leonard-atorough/tsdb/internal/models"
)

type Reader struct {
	filePath string
}

func NewReader(filePath string) *Reader {
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	fullPath := filepath.Join(projectRoot, filePath)
	return &Reader{
		filePath: fullPath,
	}
}

func (r *Reader) Query(startTime, endTime string) ([]models.TimeSeriesData, error) {
	startTimeMilli, err := convertTimeToUnix(startTime)
	if err != nil {
		return nil, err
	}
	endTimeMilli, err := convertTimeToUnix(endTime)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []models.TimeSeriesData
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data, err := models.UnmarshalLine(scanner.Bytes())
		if err != nil {
			return nil, err
		}
		if data.Timestamp >= startTimeMilli && data.Timestamp <= endTimeMilli {
			results = append(results, *data)
		}
		if data.Timestamp > endTimeMilli {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}