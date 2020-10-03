package ukishima

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
	s := NewScrapingServiceUkishimaNishiizu("ukishima-boat-in-shizuoka-nishiizu", "http://srdkaikyo.sblo.jp/", logging)
	pwd, _ := os.Getwd()
	s.url = pwd + "/testdata/" + t.Name() + "_20201003.html"
	ocean, _ := s.Scrape()
	assert.Equal(t, &model.Ocean{
		LocationName: "ukishima-boat-in-shizuoka-nishiizu",
		URL:          s.url,
		Temperature: model.Temperature{
			Min: 24,
			Med: -1,
			Max: 26,
		},
		Visibility: model.Visibility{
			Min: -1,
			Med: 20,
			Max: -1,
		},
		MeasuredTime: time.Date(2020, time.October, 3, 0, 0, 0, 0, time.Local),
	}, ocean)
}
