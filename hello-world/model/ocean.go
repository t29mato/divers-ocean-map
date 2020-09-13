package model

import "time"

// Ocean ...
type Ocean struct {
	temperature Temperature
	visibility  Visibility
	createdAt   time.Time
	measuredAt  time.Time
}

// Temperature ...
type Temperature struct {
	avg int
	min int
	max int
}

// Visibility ...
type Visibility struct {
	avg int
	min int
	max int
}
