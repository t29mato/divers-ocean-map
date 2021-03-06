package main

import (
	"scraping/logging"
	"scraping/repository"
	"scraping/service/shizuoka/iop"
	"scraping/service/shizuoka/ukishima"
	"scraping/service/tiba"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler ...
func Handler(e *events.CloudWatchEvent) {
	// FIXME: このイベントIDはCloudwatch Logsに出力されるIDと異なる orz
	logging := logging.NewOceanLoggingImpl(e.ID)
	logging.Info("スクレイピング開始")

	// TODO: goroutine 使う
	scrapingServiceIzuOceanPark := iop.NewFetchService("izu-ocean-park", "https://iop-dc.com/", logging)
	oceanIzuOceanPark, err := scrapingServiceIzuOceanPark.Fetch()
	if err != nil {
		logging.Info("伊豆海洋公園のスクレイピングの途中で失敗しました", err.Error())
	}

	oceanRepository := repository.NewOceanRepository(logging)

	err = oceanRepository.Store(oceanIzuOceanPark)
	if err != nil {
		logging.Info("伊豆海洋公園のDBへの挿入で失敗", err.Error())
	}

	scrapingServiceUkishimaTiba := tiba.NewFetchService("ukishima-in-tiba-katsuyama", "http://paroparo.jp", logging)
	oceanUkishimaTiba, err := scrapingServiceUkishimaTiba.Fetch()
	if err != nil {
		logging.Info("浮島 (千葉県勝山市)のスクレイピングの途中で失敗しました", err.Error())
	}

	err = oceanRepository.Store(oceanUkishimaTiba)
	if err != nil {
		logging.Info("浮島 (千葉県勝山市)のDBへの挿入で失敗", err.Error())
	}

	scrapingServiceUkishimaNishiizu := ukishima.NewFetchService("ukishima-boat-in-shizuoka-nishiizu", "http://srdkaikyo.sblo.jp/", logging)
	oceanUkishimaNishiizu, err := scrapingServiceUkishimaNishiizu.Fetch()
	if err != nil {
		logging.Info("浮島 (静岡県西伊豆)のスクレイピングの途中で失敗しました", err.Error())
	}

	err = oceanRepository.Store(oceanUkishimaNishiizu)
	if err != nil {
		logging.Info("浮島 (静岡県西伊豆)のDBへの挿入で失敗", err.Error())
	}

	logging.Info("スクレイピング終了")
}

func main() {
	lambda.Start(Handler)
}
