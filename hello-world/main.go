package main

import (
	"fmt"
	"hello-world/service"
)

var scrapingService = service.NewScrapingService()

func main() {
	url := "https://iop-dc.com/"
	ocean, err := scrapingService.Scrape(url)
	fmt.Println(ocean)
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}
}
