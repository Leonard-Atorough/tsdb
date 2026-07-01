package models

// timeSeriesData is a struct that represents a single data point in a time series database.
// It contains the measurement name, a set of tags, a set of fields, and a timestamp.
//
type TimeSeriesData struct {
	Measurement string `json:"m"`
	TagSet      map[string]string `json:"tags,omitempty"`
	FieldSet    map[string]any `json:"fields"`
	Timestamp   int64 `json:"ts"`
}