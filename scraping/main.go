package main

import (
	"scraping/logging"
	"scraping/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler ...
func Handler(e *events.CloudWatchEvent) {
	// FIXME: このイベントIDはCloudwatch Logsに出力されるIDと異なる orz
	logging := logging.NewOceanLoggingImpl(e.ID)
	logging.Info("スクレイピング開始")

	// TODO: goroutine 使う
	scrapingServiceIzuOceanPark := service.NewScrapingServiceIzuOceanPark(logging)
	oceanIzuOceanPark, err := scrapingServiceIzuOceanPark.Scrape()
	if err != nil {
		logging.Info("伊豆海洋公園のスクレイピングの途中で失敗しました", err.Error())
	}

	err = scrapingServiceIzuOceanPark.ScrapingService.Store(oceanIzuOceanPark)
	if err != nil {
		logging.Info("伊豆海洋公園のDBへの挿入で失敗", err.Error())
	}

	scrapingServiceUkishimaTiba := service.NewScrapingServiceUkishimaTiba(logging)
	oceanUkishimaTiba, err := scrapingServiceUkishimaTiba.Scrape()
	if err != nil {
		logging.Info("浮島 (千葉県勝山市)のスクレイピングの途中で失敗しました", err.Error())
	}

	err = scrapingServiceUkishimaTiba.ScrapingService.Store(oceanUkishimaTiba)
	if err != nil {
		logging.Info("浮島 (千葉県勝山市)のDBへの挿入で失敗", err.Error())
	}

	logging.Info("スクレイピング終了")
}

func main() {
	lambda.Start(Handler)
}
