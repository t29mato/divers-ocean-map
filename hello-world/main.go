package main

import (
	"fmt"
	"hello-world/service"
)

var scrapingService = service.NewScrapingService()
var dynamodbService = service.NewDynamoDBService()

func main() {
	url := "https://iop-dc.com/"
	ocean, err := scrapingService.Scrape(url)
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}

	err = dynamodbService.Create(ocean)
	if err != nil {
		fmt.Println("DBへの挿入で失敗", err)
	}

}
