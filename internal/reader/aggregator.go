package reader

import "github.com/leonard-atorough/tsdb/internal/models"

func (r *Reader) Avg(measurement, field string, startTime, endTime string) (float64, error) {
	points, err := r.Query(startTime, endTime)
	if err != nil {
		return 0, err
	}

	return sumPoints(points, measurement, field) / float64(countPoints(points, measurement, field)), nil
}

func (r *Reader) Sum(measurement, field string, startTime, endTime string) (float64, error) {
	points, err := r.Query(startTime, endTime)
	if err != nil {
		return 0, err
	}

	return sumPoints(points, measurement, field), nil
}

func (r *Reader) Count(measurement, field string, startTime, endTime string) (int, error) {
	points, err := r.Query(startTime, endTime)
	if err != nil {
		return 0, err
	}

	return countPoints(points, measurement, field), nil
}

func (r *Reader) Min(measurement, field string, startTime, endTime string) (float64, error) {
	points, err := r.Query(startTime, endTime)
	if err != nil {
		return 0, err
	}

	var min float64
	first := true
	for _, point := range points {
		if matchesMeasurementAndField(point, measurement, field) {
			if first {
				min = point.FieldSet[field].(float64)
				first = false
			} else if point.FieldSet[field].(float64) < min {
				min = point.FieldSet[field].(float64)
			}
		}
	}
	return min, nil
}

func (r *Reader) Max(measurement, field string, startTime, endTime string) (float64, error) {
	points, err := r.Query(startTime, endTime)
	if err != nil {
		return 0, err
	}

	var max float64
	first := true
	for _, point := range points {
		if matchesMeasurementAndField(point, measurement, field) {
			if first {
				max = point.FieldSet[field].(float64)
				first = false
			} else if point.FieldSet[field].(float64) > max {
				max = point.FieldSet[field].(float64)
			}
		}
	}
	return max, nil
}

func (r *Reader) Aggregates(measurement, field string, funcs []string, startTime, endTime string) (map[string]float64, error) {
	results := make(map[string]float64)

	for _, fn := range funcs {
		switch fn {
		case "sum":
			results["sum"], _ = r.Sum(measurement, field, startTime, endTime)
		case "count":
			count, _ := r.Count(measurement, field, startTime, endTime)
			results["count"] = float64(count)
		case "avg":
			results["avg"], _ = r.Avg(measurement, field, startTime, endTime)
		case "min":
			results["min"], _ = r.Min(measurement, field, startTime, endTime)
		case "max":
			results["max"], _ = r.Max(measurement, field, startTime, endTime)
		}
	}

	return results, nil
}

func sumPoints(points []models.TimeSeriesData, measurement, field string) float64 {
	var sum float64
	for _, point := range points {
		if matchesMeasurementAndField(point, measurement, field) {
			sum += point.FieldSet[field].(float64)
		}
	}
	return sum
}

func countPoints(points []models.TimeSeriesData, measurement, field string) int {
	var count int
	for _, point := range points {
		if matchesMeasurementAndField(point, measurement, field) {
			count++
		}
	}
	return count
}

func matchesMeasurementAndField(point models.TimeSeriesData, measurement, field string) bool {
	if point.Measurement != measurement {
		return false
	}
	_, ok := point.FieldSet[field]
	return ok
}
