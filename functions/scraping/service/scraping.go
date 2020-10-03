package service

import (
	"scraping/model"
)

// ScrapingService ...
type ScrapingService interface {
	Scrape() (*model.Ocean, error)
}
