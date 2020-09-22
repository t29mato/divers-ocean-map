package main

import (
	"fmt"
	"hello-world/service"
)

func main() {
	// TODO: goroutine 使う
	scrapingServiceIzuOceanPark := service.NewScrapingServiceIzuOceanPark()
	oceanIzuOceanPark, err := scrapingServiceIzuOceanPark.Scrape()
	if err != nil {
		fmt.Println("伊豆海洋公園のスクレイピングの途中で失敗しました", err)
	}

	err = scrapingServiceIzuOceanPark.ScrapingService.Store(oceanIzuOceanPark)
	if err != nil {
		fmt.Println("伊豆海洋公園のDBへの挿入で失敗", err)
	}

	scrapingServiceUkishimaTiba := service.NewScrapingServiceUkishimaTiba()
	oceanUkishimaTiba, err := scrapingServiceUkishimaTiba.Scrape()
	if err != nil {
		fmt.Println("浮島 (千葉県勝山市)のスクレイピングの途中で失敗しました", err)
	}

	err = scrapingServiceUkishimaTiba.ScrapingService.Store(oceanUkishimaTiba)
	if err != nil {
		fmt.Println("浮島 (千葉県勝山市)のDBへの挿入で失敗", err)
	}
}
