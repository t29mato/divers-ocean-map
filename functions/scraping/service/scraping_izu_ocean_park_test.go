package service

import (
	"os"
	"scraping/logging"
	"scraping/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestScrape ...
func TestScrapeIOP(t *testing.T) {
	logging := logging.NewOceanLoggingImpl("66936b3e-08e3-404b-815d-ddbccfb03cc9")
	s := NewScrapingServiceIzuOceanPark(logging)
	pwd, _ := os.Getwd()
	url := pwd + "/testdata/" + t.Name() + "_20200913.html"
	s.ScrapingService.url = url
	ocean, _ := s.Scrape()
	assert.Equal(t, &model.Ocean{
		LocationName: "伊豆海洋公園",
		Temperature: model.Temperature{
			Min: 21,
			Med: -1,
			Max: 27,
		},
		Visibility: model.Visibility{
			Min: 10,
			Med: -1,
			Max: 20,
		},
		MeasuredTime: time.Date(2020, time.September, 13, 0, 0, 0, 0, time.UTC),
	}, ocean)

	url = pwd + "/testdata/" + t.Name() + "_20200922.html"
	s.ScrapingService.url = url
	ocean, _ = s.Scrape()
	assert.Equal(t, &model.Ocean{
		LocationName: "伊豆海洋公園",
		Temperature: model.Temperature{
			Min: 23,
			Med: -1,
			Max: 26,
		},
		Visibility: model.Visibility{
			Min: -1,
			Med: 15,
			Max: -1,
		},
		MeasuredTime: time.Date(2020, time.September, 22, 0, 0, 0, 0, time.UTC),
	}, ocean)
}
