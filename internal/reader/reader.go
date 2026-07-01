package reader

import (
	"bufio"
	"os"
	"path/filepath"
	"time"

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
	From        string
	To          string
	Ago         string
	Measurement string
	Tags        map[string]string
}

func (r *Reader) Query(opts QueryOpts) ([]models.TimeSeriesData, error) {
	fromMs, toMs, err := r.resolveTimeBounds(opts.From, opts.To, opts.Ago)
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
			// Skip invalid or unparsable lines rather than failing the whole query.
			continue
		}
		if data.Timestamp > toMs {
			break
		}

		if !matchesMeasurement(data.Measurement, opts.Measurement) {
			continue
		}

		if !matchesTags(data.TagSet, opts.Tags) {
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

func (r *Reader) resolveTimeBounds(from, to, ago string) (int64, int64, error) {
	now := time.Now()

	// Resolve "To"
	if to == "" {
		toMs := now.UnixMilli()

		// Resolve "From"
		switch {
		case from != "":
			t, err := time.Parse(time.RFC3339, from)
			if err != nil {
				return 0, 0, err
			}
			fromMs := t.UnixMilli()
			return fromMs, toMs, nil

		case ago != "":
			durationNs, err := time.ParseDuration(ago)
			if err != nil {
				return 0, 0, err
			}
			fromMs := toMs - int64(durationNs.Milliseconds())
			return fromMs, toMs, nil

		default:
			// Default 1 week
			fromMs := toMs - (7 * 24 * 60 * 60 * 1000)
			return fromMs, toMs, nil
		}
	}

	// Both from and to explicitly provided
	fromTs, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return 0, 0, err
	}
	toTs, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return 0, 0, err
	}

	return fromTs.UnixMilli(), toTs.UnixMilli(), nil
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

// matchesMeasurement checks if the measurement in the data matches the query measurement.
func matchesMeasurement(dataMeasurement string, queryMeasurement string) bool {
	// If no measurement filter is provided, accept any measurement.
	if queryMeasurement == "" {
		return true
	}
	return dataMeasurement == queryMeasurement
}

// matchesTimeRange checks if the timestamp of the data point falls within the specified time range.
func matchesTimeRange(timestamp int64, fromMs int64, toMs int64) bool {
	return timestamp >= fromMs && timestamp <= toMs
}
