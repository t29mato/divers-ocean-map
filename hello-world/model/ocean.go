package model

import "time"

// Ocean ...
type Ocean struct {
	Temperature  Temperature
	Visibility   Visibility
	MeasuredTime time.Time
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
