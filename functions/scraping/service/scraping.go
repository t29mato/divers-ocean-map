package service

import (
	"scraping/logging"
	"scraping/model"
)

// ScrapingService ...
type ScrapingService interface {
	Scrape() (*model.Ocean, error)
}

// ScrapingServiceImpl ...
type ScrapingServiceImpl struct {
	logging *logging.OceanLoggingImpl
	db      *DynamoDBServiceImpl
}

// Store ...
func (s *ScrapingServiceImpl) Store(ocean *model.Ocean) error {
	s.logging.Info("データの永久保存開始", ocean.LocationName)
	err := s.db.CreateIfNotExist(ocean)
	if err != nil {
		s.logging.Info("データの永久保存に失敗")
		return err
	}
	s.logging.Info("データの永久保存終了", ocean.LocationName)
	return nil
}
