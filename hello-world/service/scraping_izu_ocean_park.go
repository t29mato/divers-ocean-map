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

// ScrapingServiceIzuOceanParkImpl ...
type ScrapingServiceIzuOceanParkImpl struct {
	ScrapingService   *ScrapingServiceImpl
	queryTemperature  string
	queryVisibility   string
	queryMeasuredTime string
}

// NewScrapingServiceIzuOceanPark ...
func NewScrapingServiceIzuOceanPark() *ScrapingServiceIzuOceanParkImpl {
	return &ScrapingServiceIzuOceanParkImpl{
		ScrapingService: &ScrapingServiceImpl{
			url: "https://iop-dc.com/",
			db:  NewDynamoDBService(),
		},
		queryTemperature:  "#homeConditionDetail > dl > dd:nth-child(2)",
		queryVisibility:   "#homeConditionDetail > dl > dd:nth-child(4)",
		queryMeasuredTime: "#homeConditionIndex > dl > dt",
	}
}

// Scrape ...
func (s *ScrapingServiceIzuOceanParkImpl) Scrape() (*model.Ocean, error) {
	ocean := model.NewOcean()

	// DOM取得
	doc, err := s.fetchDocument(s.ScrapingService.url)
	if err != nil {
		fmt.Println("HTMLファイルの読み込みに失敗しました。url =", s.ScrapingService.url)
		return nil, err
	}

	// 水温取得
	err = s.fetchTemperature(s.queryTemperature, doc, ocean)
	if err != nil {
		fmt.Println("水温の取得に失敗")
		return nil, err
	}

	// 透明度取得
	err = s.fetchVisibility(s.queryVisibility, doc, ocean)
	if err != nil {
		fmt.Println("透明度の取得に失敗")
		return nil, err
	}

	// 測定日時取得
	err = s.fetchMeasuredTime(s.queryMeasuredTime, doc, ocean)
	if err != nil {
		fmt.Println("測定日時の取得に失敗")
		return nil, err
	}

	return ocean, err
}

func (s *ScrapingServiceIzuOceanParkImpl) fetchDocument(url string) (*goquery.Document, error) {
	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if strings.Contains(url, "http") {
		return goquery.NewDocument(url)
	}
	file, err := os.Open(url)
	if err != nil {
		fmt.Println("有効なファイルパスでありません。url =", url)
		return nil, err
	}
	defer file.Close()
	return goquery.NewDocumentFromReader(file)
}

func (s *ScrapingServiceIzuOceanParkImpl) fetchTemperature(query string, doc *goquery.Document, ocean *model.Ocean) error {
	temperatureHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	temperatures := reg.FindAllStringSubmatch(temperatureHTML, -1)
	min, err := strconv.Atoi(temperatures[0][0])
	if err != nil {
		fmt.Println("水温1の数値変換に失敗, 変換前=", temperatures[0][0])
		return err
	}
	switch len(temperatures) {
	case 1:
		ocean.Temperature.Med = min
	case 2:
		max, err := strconv.Atoi(temperatures[1][0])
		if err != nil {
			fmt.Println("水温2の数値変換に失敗, 変換前=", temperatures[1][0])
			return err
		}
		ocean.Temperature.Min = min
		ocean.Temperature.Max = max
	}
	return nil
}

func (s *ScrapingServiceIzuOceanParkImpl) fetchVisibility(query string, doc *goquery.Document, ocean *model.Ocean) error {
	visibilityHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	visibilities := reg.FindAllStringSubmatch(visibilityHTML, -1)

	min, err := strconv.Atoi(visibilities[0][0])
	if err != nil {
		fmt.Println("透明度1の数値変換に失敗, 変換前=", visibilities[0][0])
		return err
	}
	switch len(visibilities) {
	case 1:
		ocean.Visibility.Med = min
	case 2:
		max, err := strconv.Atoi(visibilities[1][0])
		if err != nil {
			fmt.Println("透明度2の数値変換に失敗, 変換前=", visibilities[1][0])
			return err
		}
		ocean.Visibility.Min = min
		ocean.Visibility.Max = max
	}
	return nil
}

func (s *ScrapingServiceIzuOceanParkImpl) fetchMeasuredTime(query string, doc *goquery.Document, ocean *model.Ocean) error {
	MeasuredTimeHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	measuredTimes := reg.FindAllStringSubmatch(MeasuredTimeHTML, -1)

	month, err := strconv.Atoi(measuredTimes[0][0])
	if err != nil {
		fmt.Println("月の数値変換に失敗, 変換前=", measuredTimes[0][0])
		return err
	}
	day, err := strconv.Atoi(measuredTimes[1][0])
	if err != nil {
		fmt.Println("日の数値変換に失敗, 変換前=", measuredTimes[1][0])
		return err
	}
	ocean.MeasuredTime = time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return nil
}
