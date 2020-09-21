package service

import (
	"fmt"
	"hello-world/model"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ss ScrapingServiceImpl

// TestScrape ...
func TestScrapeIOP(t *testing.T) {
	pwd, _ := os.Getwd()
	url := pwd + "/testdata/" + t.Name() + "_20200913.html"
	fmt.Println(url)
	ocean, _ := ss.Scrape(url)
	assert.Equal(t, &model.Ocean{
		Name: "",
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
