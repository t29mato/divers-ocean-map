package main

import (
	"fmt"
	"hello-world/service"
)

func main() {
	scrapingService := service.NewScrapingServiceIzuOceanPark()
	ocean, err := scrapingService.Scrape()
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}

	err = scrapingService.ScrapingService.Store(ocean)
	if err != nil {
		fmt.Println("DBへの挿入で失敗", err)
	}

}
