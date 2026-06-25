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

// TimeSeriesDataQuery is a struct that represents a query for time series data.
// It contains a collection of TimeSeriesData and a query string.
// Not currently used, but could be useful for future implementations.
type TimeSeriesDataQuery struct {
	Collection []TimeSeriesData
	Query      string
}