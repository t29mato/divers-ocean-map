package service

import (
	"scraping/model"

	"github.com/PuerkitoBio/goquery"
)

// FetchService ...
type FetchService interface {
	Fetch() (*model.Ocean, error)
	fetchDocument(url string, ocean *model.Ocean) (*goquery.Document, error)
	fetchTemperature(query string, doc *goquery.Document, ocean *model.Ocean) error
	fetchVisibility(query string, doc *goquery.Document, ocean *model.Ocean) error
	fetchMeasuredTime(query string, doc *goquery.Document, ocean *model.Ocean) error
}
