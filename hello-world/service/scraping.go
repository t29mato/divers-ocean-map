package service

import (
	"hello-world/model"
)

// ScrapingService ...
type ScrapingService interface {
	Scrape() (*model.Ocean, error)
}

// ScrapingServiceImpl ...
type ScrapingServiceImpl struct {
	url string
	db  *DynamoDBServiceImpl
}

// Store ...
func (s *ScrapingServiceImpl) Store(ocean *model.Ocean) error {
	err := s.db.Create(ocean)
	if err != nil {
		return err
	}
	return nil
}
