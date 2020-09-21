package main

import (
	"fmt"
	"hello-world/service"
)

var scrapingService = service.NewScrapingServiceIzuOceanPark()
var dynamodbService = service.NewDynamoDBService()

func main() {
	ocean, err := scrapingService.Scrape()
	if err != nil {
		fmt.Println("スクレイピングの途中で失敗しました", err)
	}

	err = dynamodbService.Create(ocean)
	if err != nil {
		fmt.Println("DBへの挿入で失敗", err)
	}

}
