package service

import (
	"fmt"
	"hello-world/model"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ScrapingServiceImpl ...
type ScrapingServiceImpl struct{}

// ScrapingService ...
type ScrapingService interface {
	Scrape(url string) (*model.Ocean, error)
}

// NewScrapingService ...
func NewScrapingService() ScrapingService {
	return &ScrapingServiceImpl{}
}

// Scrape ...
func (s *ScrapingServiceImpl) Scrape(url string) (*model.Ocean, error) {
	var doc *goquery.Document
	var err error
	var ocean model.Ocean

	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if strings.Contains(url, "http") {
		doc, err = goquery.NewDocument(url)
	} else {
		file, err := os.Open(url)
		if err != nil {
			fmt.Println("有効なファイルパスでありません。url =", url)
			return nil, err
		}
		defer file.Close()
		doc, err = goquery.NewDocumentFromReader(file)
	}

	temperatureHTML, _ := doc.Find("#homeConditionDetail > dl > dd:nth-child(2)").Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	temperatures := reg.FindAllStringSubmatch(temperatureHTML, -1)

	min, err := strconv.Atoi(temperatures[0][0])
	if err != nil {
		fmt.Println("水温1の数値変換に失敗, 変換前=", temperatures[0][0])
		return nil, err
	}
	switch len(temperatures) {
	case 1:
		ocean.Temperature.Med = min
	case 2:
		max, err := strconv.Atoi(temperatures[1][0])
		if err != nil {
			fmt.Println("水温2の数値変換に失敗, 変換前=", temperatures[1][0])
			return nil, err
		}
		ocean.Temperature.Min = min
		ocean.Temperature.Max = max
	}

	visibilityHTML, _ := doc.Find("#homeConditionDetail > dl > dd:nth-child(4)").Html()
	visibilities := reg.FindAllStringSubmatch(visibilityHTML, -1)

	min, err = strconv.Atoi(visibilities[0][0])
	if err != nil {
		fmt.Println("透明度1の数値変換に失敗, 変換前=", visibilities[0][0])
		return nil, err
	}
	switch len(visibilities) {
	case 1:
		ocean.Visibility.Med = min
	case 2:
		max, err := strconv.Atoi(visibilities[1][0])
		if err != nil {
			fmt.Println("透明度2の数値変換に失敗, 変換前=", visibilities[1][0])
			return nil, err
		}
		ocean.Visibility.Min = min
		ocean.Visibility.Max = max
	}

	MeasuredTimeHTML, _ := doc.Find("#homeConditionIndex > dl > dt").Html()
	measuredTimes := reg.FindAllStringSubmatch(MeasuredTimeHTML, -1)

	month, err := strconv.Atoi(measuredTimes[0][0])
	if err != nil {
		fmt.Println("月の数値変換に失敗, 変換前=", measuredTimes[0][0])
		return nil, err
	}
	day, err := strconv.Atoi(measuredTimes[1][0])
	if err != nil {
		fmt.Println("日の数値変換に失敗, 変換前=", measuredTimes[1][0])
		return nil, err
	}
	ocean.MeasuredTime = time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return &ocean, err
}
