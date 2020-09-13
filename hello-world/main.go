package main

import (
	"fmt"
	"hello-world/service"
)

var scrapingService service.ScrapingService

func main() {
	url := "https://iop-dc.com/"
	err := scrapingService.Scrape(url)
	if err != nil {
		fmt.Println("url scrapping failed")
	}
}
