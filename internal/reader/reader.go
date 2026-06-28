package reader

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/leonard-atorough/tsdb/internal"
	"github.com/leonard-atorough/tsdb/internal/models"
)

type Reader struct {
	filePath string
}

func NewReader(filePath string) (*Reader, error) {
	if filepath.IsAbs(filePath) {
		return &Reader{
			filePath: filePath,
		}, nil
	} else {
		projectRoot, err := internal.GetProjectRoot()
		if err != nil {
			return nil, err
		}
		fullPath := filepath.Join(projectRoot, filePath)
		return &Reader{
			filePath: fullPath,
		}, nil
	}
}

type QueryOpts struct {
	from        string
	to          string
	Measurement string
	tags        map[string]string
}

func (r *Reader) Query(opts QueryOpts) ([]models.TimeSeriesData, error) {
	fromMs, err := internal.ConvertTimeToUnix(opts.from)
	if err != nil {
		return nil, err
	}
	if opts.to == "" {
		opts.to = internal.GetCurrentTime()
	}
	toMs, err := internal.ConvertTimeToUnix(opts.to)
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
		if data.Timestamp > toMs {
			break
		}

		if !matchesMeasurement(data.Measurement, opts.Measurement) {
			continue
		}

		if !matchesTags(data.TagSet, opts.tags) {
			continue
		}

		if !matchesTimeRange(data.Timestamp, fromMs, toMs) {
			continue
		}

		results = append(results, *data)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// matchesTags checks if the tags in the data match the query tags.
func matchesTags(dataTags map[string]string, queryTags map[string]string) bool {
	for key, value := range queryTags {
		if dataValue, exists := dataTags[key]; !exists || dataValue != value {
			return false
		}
	}
	return true
}

func matchesMeasurement(dataMeasurement string, queryMeasurement string) bool {
	return dataMeasurement == queryMeasurement
}

func matchesTimeRange(timestamp int64, fromMs int64, toMs int64) bool {
	return timestamp >= fromMs && timestamp <= toMs
}
