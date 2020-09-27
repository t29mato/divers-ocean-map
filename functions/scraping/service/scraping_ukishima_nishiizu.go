package service

import (
	"fmt"
	"os"
	"regexp"
	"scraping/logging"
	"scraping/model"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ScrapingServiceUkishimaNishiizuImpl ...
type ScrapingServiceUkishimaNishiizuImpl struct {
	ScrapingService *ScrapingServiceImpl
	queryArticle    string
	queryDate       string
}

// NewScrapingServiceUkishimaNishiizu ...
func NewScrapingServiceUkishimaNishiizu(logging *logging.OceanLoggingImpl) *ScrapingServiceUkishimaNishiizuImpl {
	// HACK: ScrapingServiceImplの中で、各ダイビングポイントの場所に依存して変わるもの以外は、ダイビングポイントの構造体に直接持たせる (隠したい)
	return &ScrapingServiceUkishimaNishiizuImpl{
		ScrapingService: &ScrapingServiceImpl{
			url:     "http://srdkaikyo.sblo.jp/",
			db:      NewDynamoDBService(),
			logging: logging,
		},
		queryArticle: "#content > div:nth-child(2) > div > div.text",
		queryDate:    "#content > div:nth-child(2) > h2",
	}
}

// Scrape ...
func (s *ScrapingServiceUkishimaNishiizuImpl) Scrape() (*model.Ocean, error) {
	s.ScrapingService.logging.Info("浮島(西伊豆)のスクレイピング開始")
	ocean := model.NewOcean("ukishima-boat-in-shizuoka-nishiizu")

	// DOM取得
	doc, err := s.fetchDocument(s.ScrapingService.url)
	if err != nil {
		s.ScrapingService.logging.Info("HTMLファイルの読み込みに失敗しました。url =", s.ScrapingService.url)
		return nil, err
	}

	// 水温取得
	err = s.fetchTemperature(s.queryArticle, doc, ocean)
	if err != nil {
		s.ScrapingService.logging.Info("水温の取得に失敗")
		return nil, err
	}

	// 透明度取得
	err = s.fetchVisibility(s.queryArticle, doc, ocean)
	if err != nil {
		s.ScrapingService.logging.Info("透明度の取得に失敗")
		return nil, err
	}

	// 測定日時取得
	err = s.fetchMeasuredTime(s.queryDate, doc, ocean)
	if err != nil {
		s.ScrapingService.logging.Info("測定日時の取得に失敗")
		return nil, err
	}

	s.ScrapingService.logging.Info("伊豆海洋公園のスクレイピング終了")
	return ocean, err
}

func (s *ScrapingServiceUkishimaNishiizuImpl) fetchDocument(url string) (*goquery.Document, error) {
	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if strings.Contains(url, "http") {
		return goquery.NewDocument(url)
	}
	file, err := os.Open(url)
	if err != nil {
		s.ScrapingService.logging.Info("有効なファイルパスでありません。url =", url)
		return nil, err
	}
	defer file.Close()
	return goquery.NewDocumentFromReader(file)
}

func (s *ScrapingServiceUkishimaNishiizuImpl) fetchTemperature(query string, doc *goquery.Document, ocean *model.Ocean) error {
	articleHTML, _ := doc.Find(query).Html()
	fmt.Println(articleHTML)
	reg := regexp.MustCompile(`水温[\s\S]*℃`)
	temperatureHTML := reg.FindAllStringSubmatch(articleHTML, -1)
	reg = regexp.MustCompile(`[0-9０-９]{1,2}`)
	temperatures := reg.FindAllStringSubmatch(temperatureHTML[0][0], -1)

	min, err := convertIntFromFullWidthString(&temperatures[0][0])
	if err != nil {
		s.ScrapingService.logging.Info("水温1の数値変換に失敗, 変換前=", temperatures[0][0])
		return err
	}
	switch len(temperatures) {
	case 1:
		ocean.Temperature.Med = min
	case 2:
		max, err := convertIntFromFullWidthString(&temperatures[1][0])
		if err != nil {
			s.ScrapingService.logging.Info("水温2の数値変換に失敗, 変換前=", temperatures[1][0])
			return err
		}
		ocean.Temperature.Min = min
		ocean.Temperature.Max = max
	}
	return nil
}

func (s *ScrapingServiceUkishimaNishiizuImpl) fetchVisibility(query string, doc *goquery.Document, ocean *model.Ocean) error {
	articleHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`透明度[\s\S]*ｍ`)
	visibilityHTML := reg.FindAllStringSubmatch(articleHTML, -1)
	reg = regexp.MustCompile(`[0-9０-９]{1,2}`)
	visibilities := reg.FindAllStringSubmatch(visibilityHTML[0][0], -1)

	min, err := convertIntFromFullWidthString(&visibilities[0][0])
	if err != nil {
		s.ScrapingService.logging.Info("透明度1の数値変換に失敗, 変換前=", visibilities[0][0])
		return err
	}
	switch len(visibilities) {
	case 1:
		ocean.Visibility.Med = min
	case 2:
		max, err := convertIntFromFullWidthString(&visibilities[1][0])
		if err != nil {
			s.ScrapingService.logging.Info("透明度2の数値変換に失敗, 変換前=", visibilities[1][0])
			return err
		}
		ocean.Visibility.Min = min
		ocean.Visibility.Max = max
	}
	return nil
}

func (s *ScrapingServiceUkishimaNishiizuImpl) fetchMeasuredTime(query string, doc *goquery.Document, ocean *model.Ocean) error {
	MeasuredTimeHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,4}`)
	measuredTimes := reg.FindAllStringSubmatch(MeasuredTimeHTML, -1)

	// HACK
	year, err := strconv.Atoi(measuredTimes[0][0])
	month, err := strconv.Atoi(measuredTimes[1][0])
	day, err := strconv.Atoi(measuredTimes[2][0])
	if err != nil {
		s.ScrapingService.logging.Info("datetimeの変換に失敗")
		return err
	}
	ocean.MeasuredTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return nil
}
