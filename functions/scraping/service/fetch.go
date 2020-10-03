package service

import (
	"scraping/model"

	"github.com/PuerkitoBio/goquery"
)

// FetchService ...
type FetchService interface {
	Fetch() (*model.Ocean, error)
	fetchDocument(url string) (*goquery.Document, error)
}
