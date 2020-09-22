package main

import (
	"fmt"
	"log"
	"scraping/logging"
	"scraping/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler ...
func Handler(e *events.CloudWatchEvent) {
	logging := logging.NewOceanLoggingImpl(e.ID)
	log.Println("スクレイピング開始")

	// TODO: goroutine 使う
	scrapingServiceIzuOceanPark := service.NewScrapingServiceIzuOceanPark(logging)
	oceanIzuOceanPark, err := scrapingServiceIzuOceanPark.Scrape()
	if err != nil {
		fmt.Println("伊豆海洋公園のスクレイピングの途中で失敗しました", err)
	}

	err = scrapingServiceIzuOceanPark.ScrapingService.Store(oceanIzuOceanPark)
	if err != nil {
		fmt.Println("伊豆海洋公園のDBへの挿入で失敗", err)
	}

	scrapingServiceUkishimaTiba := service.NewScrapingServiceUkishimaTiba(logging)
	oceanUkishimaTiba, err := scrapingServiceUkishimaTiba.Scrape()
	if err != nil {
		fmt.Println("浮島 (千葉県勝山市)のスクレイピングの途中で失敗しました", err)
	}

	err = scrapingServiceUkishimaTiba.ScrapingService.Store(oceanUkishimaTiba)
	if err != nil {
		fmt.Println("浮島 (千葉県勝山市)のDBへの挿入で失敗", err)
	}

	fmt.Println("スクレイピング終了")
}

func main() {
	lambda.Start(Handler)
}
