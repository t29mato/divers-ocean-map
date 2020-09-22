package main

import (
	"fmt"
	"hello-world/service"
)

func main() {
	scrapingServiceIzuOceanPark := service.NewScrapingServiceIzuOceanPark()
	oceanIzuOceanPark, err := scrapingServiceIzuOceanPark.Scrape()
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}

	err = scrapingServiceIzuOceanPark.ScrapingService.Store(oceanIzuOceanPark)
	if err != nil {
		fmt.Println("DBへの挿入で失敗", err)
	}

	scrapingServiceUkishimaTiba := service.NewScrapingServiceUkishimaTiba()
	oceanUkishimaTiba, err := scrapingServiceUkishimaTiba.Scrape()
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}

	err = scrapingServiceUkishimaTiba.ScrapingService.Store(oceanUkishimaTiba)
	if err != nil {
		fmt.Println("DBへの挿入で失敗", err)
	}
}
