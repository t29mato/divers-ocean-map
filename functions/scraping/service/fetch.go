package service

import (
	"scraping/model"
)

// FetchService ...
type FetchService interface {
	Fetch() (*model.Ocean, error)
}
