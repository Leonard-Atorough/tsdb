package models

import "encoding/json"

func MarshalLine(d *TimeSeriesData) ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	b = append(b, '\n')
	return b, nil
}

func UnmarshalLine(line []byte) (*TimeSeriesData, error) {
	var d TimeSeriesData
	err := json.Unmarshal(line, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}