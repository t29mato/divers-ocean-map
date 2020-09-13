package service

import (
	"fmt"
	"hello-world/model"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// TestScrape ...
func TestScrape(url string) (model.Ocean, error) {
	var doc *goquery.Document
	var err error

	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if os.Getenv("UNIT_TEST") != "true" {
		doc, err = goquery.NewDocument(url)
	} else {
		file, err := os.Open(url)
		if err != nil {
			fmt.Println("有効なファイルパスでありません。url =", url)
		}
		defer file.Close()
		doc, err = goquery.NewDocumentFromReader(file)
	}

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		fmt.Println(url)
	})

	return err
}
