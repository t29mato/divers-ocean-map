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
func TestScrapeUkishimaTiba(t *testing.T) {
	logging := logging.NewOceanLoggingImpl("66936b3e-08e3-404b-815d-ddbccfb03cc9")
	s := NewScrapingServiceUkishimaTiba(logging)
	pwd, _ := os.Getwd()
	url := pwd + "/testdata/" + t.Name() + "_20200921.html"
	s.ScrapingService.url = url
	ocean, _ := s.Scrape()
	assert.Equal(t, &model.Ocean{
		LocationName: "浮島 (千葉県勝山市)",
		Temperature: model.Temperature{
			Min: 24,
			Med: -1,
			Max: 26,
		},
		Visibility: model.Visibility{
			Min: 8,
			Med: -1,
			Max: 10,
		},
		MeasuredTime: time.Date(2020, time.September, 21, 14, 49, 6, 0, time.Local),
	}, ocean)
}
