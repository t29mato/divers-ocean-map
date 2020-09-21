package model

import "time"

// Ocean ...
type Ocean struct {
	Name         string `dynamodbav:"name"`
	Temperature  Temperature
	Visibility   Visibility
	MeasuredTime time.Time `dynamodbav:"measured_time"`
}

// NewOcean ...
func NewOcean() *Ocean {
	return &Ocean{
		Name: "",
		Temperature: Temperature{
			Min: -1,
			Med: -1,
			Max: -1,
		},
		Visibility: Visibility{
			Min: -1,
			Med: -1,
			Max: -1,
		},
		MeasuredTime: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

// Temperature ...
type Temperature struct {
	Min int `dynamodbav:"temperature_min"`
	Med int `dynamodbav:"temperature_med"`
	Max int `dynamodbav:"temperature_max"`
}

// Visibility ...
type Visibility struct {
	Min int `dynamodbav:"visibility_min"`
	Med int `dynamodbav:"visibility_med"`
	Max int `dynamodbav:"visibility_max"`
}
