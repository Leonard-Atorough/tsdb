package reader

import (
	"errors"

	"github.com/leonard-atorough/tsdb/internal/models"
)

type AggregateOpts struct {
	Field       string
	Measurement string
	Funcs       []string
	From        string
	To          string
	Tags        map[string]string
}

func (r *Reader) Avg(opts AggregateOpts) (float64, error) {
	points, err := r.Query(QueryOpts{
		From:        opts.From,
		To:          opts.To,
		Measurement: opts.Measurement,
		Tags:        opts.Tags,
	})
	if err != nil {
		return 0, err
	}

	return sumPoints(points, opts.Field) / float64(countPoints(points, opts.Field)), nil
}

func (r *Reader) Sum(opts AggregateOpts) (float64, error) {
	points, err := r.Query(QueryOpts{
		From:        opts.From,
		To:          opts.To,
		Measurement: opts.Measurement,
		Tags:        opts.Tags,
	})
	if err != nil {
		return 0, err
	}

	return sumPoints(points, opts.Field), nil
}

func (r *Reader) Count(opts AggregateOpts) (int, error) {
	points, err := r.Query(QueryOpts{
		From:        opts.From,
		To:          opts.To,
		Measurement: opts.Measurement,
		Tags:        opts.Tags,
	})
	if err != nil {
		return 0, err
	}

	return countPoints(points, opts.Field), nil
}

func (r *Reader) Min(opts AggregateOpts) (float64, error) {
	points, err := r.Query(QueryOpts{
		From:        opts.From,
		To:          opts.To,
		Measurement: opts.Measurement,
		Tags:        opts.Tags,
	})
	if err != nil {
		return 0, err
	}

	var min float64
	first := true
	for _, point := range points {
		if point.FieldSet[opts.Field] != nil {
			if first {
				min = point.FieldSet[opts.Field].(float64)
				first = false
			} else if point.FieldSet[opts.Field].(float64) < min {
				min = point.FieldSet[opts.Field].(float64)
			}
		}
	}
	return min, nil
}

func (r *Reader) Max(opts AggregateOpts) (float64, error) {
	points, err := r.Query(QueryOpts{
		From:        opts.From,
		To:          opts.To,
		Measurement: opts.Measurement,
		Tags:        opts.Tags,
	})
	if err != nil {
		return 0, err
	}

	var max float64
	first := true
	for _, point := range points {
		if point.FieldSet[opts.Field] != nil {
			if first {
				max = point.FieldSet[opts.Field].(float64)
				first = false
			} else if point.FieldSet[opts.Field].(float64) > max {
				max = point.FieldSet[opts.Field].(float64)
			}
		}
	}
	return max, nil
}

func (r *Reader) Aggregates(opts AggregateOpts) (map[string]float64, error) {
	results := make(map[string]float64)
	var errs []error

	for _, fn := range opts.Funcs {
		switch fn {
		case "sum":
			val, err := r.Sum(opts)
			if err != nil {
				errs = append(errs, err)
			} else {
				results["sum"] = val
			}
		case "count":
			count, err := r.Count(opts)
			if err != nil {
				errs = append(errs, err)
			} else {
				results["count"] = float64(count)
			}
		case "avg":
			val, err := r.Avg(opts)
			if err != nil {
				errs = append(errs, err)
			} else {
				results["avg"] = val
			}
		case "min":
			val, err := r.Min(opts)
			if err != nil {
				errs = append(errs, err)
			} else {
				results["min"] = val
			}
		case "max":
			val, err := r.Max(opts)
			if err != nil {
				errs = append(errs, err)
			} else {
				results["max"] = val
			}
		}
	}

	return results, errors.Join(errs...)
}

func sumPoints(points []models.TimeSeriesData, field string) float64 {
	var sum float64
	for _, point := range points {
		if point.FieldSet[field] != nil {
			sum += point.FieldSet[field].(float64)
		}
	}
	return sum
}

func countPoints(points []models.TimeSeriesData, field string) int {
	var count int
	for _, point := range points {
		if point.FieldSet[field] != nil {
			count++
		}
	}
	return count
}
