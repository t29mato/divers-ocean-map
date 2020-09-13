package main

import (
	"path/filepath"
	"testing"
)

func TestScrapeIOP(t *testing.T) {
	url := filepath.Join("testdata", t.Name()+".html")
	scrape(url)

}
