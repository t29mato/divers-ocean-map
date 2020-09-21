package service

import (
	"hello-world/model"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestScrape ...
func TestScrapeUkishimaTiba(t *testing.T) {
	s := NewScrapingServiceUkishimaTiba()
	pwd, _ := os.Getwd()
	url := pwd + "/testdata/" + t.Name() + "_20200921.json"
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
}
