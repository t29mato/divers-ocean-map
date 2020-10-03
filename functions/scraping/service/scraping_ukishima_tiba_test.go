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
		LocationName: "ukishima-in-tiba-katsuyama",
		URL:          "http://paroparo.jp",
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

func TestScrapeUkishimaTiba_20201003(t *testing.T) {
	logging := logging.NewOceanLoggingImpl("66936b3e-08e3-404b-815d-ddbccfb03cc9")
	s := NewScrapingServiceUkishimaTiba(logging)
	pwd, _ := os.Getwd()
	url := pwd + "/testdata/" + t.Name() + ".html"
	s.ScrapingService.url = url
	ocean, _ := s.Scrape()
	assert.Equal(t, &model.Ocean{
		LocationName: "ukishima-in-tiba-katsuyama",
		URL:          "http://paroparo.jp",
		Temperature: model.Temperature{
			Min: -1,
			Med: 24,
			Max: -1,
		},
		Visibility: model.Visibility{
			Min: 10,
			Med: -1,
			Max: 13,
		},
		MeasuredTime: time.Date(2020, time.October, 3, 16, 4, 33, 0, time.Local),
	}, ocean)
}
