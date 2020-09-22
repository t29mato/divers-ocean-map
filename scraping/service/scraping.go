package service

import (
	"fmt"
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
	err := s.db.CreateIfNotExist(ocean)
	if err != nil {
		fmt.Println("データの永久保存に失敗")
		return err
	}
	return nil
}

// convertIntFromFullWidthString
func convertIntFromFullWidthString(s *string) (int, error) {
	return strconv.Atoi(width.Narrow.String(*s))
}
