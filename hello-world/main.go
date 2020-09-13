package main

import (
	"fmt"
	"hello-world/service"
)

var scrapingService service.ScrapingService

func main() {
	url := "https://iop-dc.com/"
	ocean, err := scrapingService.Scrape(url)
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}
}
