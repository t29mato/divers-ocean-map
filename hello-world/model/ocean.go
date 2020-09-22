package model

import "time"

// Ocean ...
type Ocean struct {
	LocationName string
	Temperature  Temperature
	Visibility   Visibility
	MeasuredTime time.Time
}

// NewOcean ...
func NewOcean(locationName string) *Ocean {
	return &Ocean{
		LocationName: locationName,
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
	Min int
	Med int
	Max int
}

// Visibility ...
type Visibility struct {
	Min int
	Med int
	Max int
}
