package service

import "github.com/PuerkitoBio/goquery"

// ScrapingServiceImpl ...
type ScrapingServiceImpl struct {
	url string
}

// ScrapingService ...
type ScrapingService interface {
	setURL(url string)
	Fetch(string) *goquery.Document
}

func (s *ScrapingServiceImpl) setURL(url string) {
	s.url = url
}
