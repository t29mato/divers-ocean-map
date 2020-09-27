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
func TestScrapeUkishimaNishiizu(t *testing.T) {
	logging := logging.NewOceanLoggingImpl("66936b3e-08e3-404b-815d-ddbccfb03cc9")
	s := NewScrapingServiceUkishimaNishiizu(logging)
	pwd, _ := os.Getwd()
	url := pwd + "/testdata/" + t.Name() + "_20200927.html"
	s.ScrapingService.url = url
	ocean, _ := s.Scrape()
	assert.Equal(t, &model.Ocean{
		LocationName: "ukishima-boat-in-shizuoka-nishiizu",
		Temperature: model.Temperature{
			Min: 25,
			Med: -1,
			Max: 29,
		},
		Visibility: model.Visibility{
			Min: 15,
			Med: -1,
			Max: 20,
		},
		MeasuredTime: time.Date(2020, time.September, 27, 0, 0, 0, 0, time.Local),
	}, ocean)
}
