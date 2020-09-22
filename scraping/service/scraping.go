package service

import (
	"scraping/logging"
	"scraping/model"
	"strconv"

	"golang.org/x/text/width"
)

// ScrapingService ...
type ScrapingService interface {
	Scrape() (*model.Ocean, error)
}

// ScrapingServiceImpl ...
type ScrapingServiceImpl struct {
	url     string
	db      *DynamoDBServiceImpl
	logging *logging.OceanLoggingImpl
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

// convertIntFromFullWidthString
func convertIntFromFullWidthString(s *string) (int, error) {
	return strconv.Atoi(width.Narrow.String(*s))
}
