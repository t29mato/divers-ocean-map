package service

import (
	"os"
	"regexp"
	"scraping/logging"
	"scraping/model"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ScrapingServiceUkishimaTibaImpl ...
type ScrapingServiceUkishimaTibaImpl struct {
	ScrapingService   *ScrapingServiceImpl
	queryArticle      string
	queryMeasuredTime string
}

// NewScrapingServiceUkishimaTiba ...
func NewScrapingServiceUkishimaTiba(logging *logging.OceanLoggingImpl) *ScrapingServiceUkishimaTibaImpl {
	return &ScrapingServiceUkishimaTibaImpl{
		ScrapingService: &ScrapingServiceImpl{
			url:     "http://paroparo.jp",
			db:      NewDynamoDBService(),
			logging: logging,
		},
		queryArticle:      "div.entry-content",
		queryMeasuredTime: "#homeConditionIndex > dl > dt",
	}
}

// Scrape ...
func (s *ScrapingServiceUkishimaTibaImpl) Scrape() (*model.Ocean, error) {
	ocean := model.NewOcean("ukishima-in-tiba-katsuyama")

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
	err = s.fetchMeasuredTime("footer > span.posted-on > a:nth-child(2) > time.entry-date.published", doc, ocean)
	if err != nil {
		s.ScrapingService.logging.Info("測定日時の取得に失敗")
		return nil, err
	}

	return ocean, err
}

func (s *ScrapingServiceUkishimaTibaImpl) fetchDocument(url string) (*goquery.Document, error) {
	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if strings.Contains(url, "http") {
		doc, _ := goquery.NewDocument(url)
		// トップページには透明度情報がないので、トップページから最新の記事のURLを取得する
		latestArticleURL, _ := doc.Find("#post-9 > div > div:nth-child(5) > div > ul > li:nth-child(1) > div.kaiyou_thumb > a").Attr("href")
		s.ScrapingService.logging.Info("latestArticleURL:", latestArticleURL)
		latestDoc, _ := goquery.NewDocument(latestArticleURL)
		return latestDoc, nil
	}
	file, err := os.Open(url)
	if err != nil {
		s.ScrapingService.logging.Info("有効なファイルパスでありません。url =", url)
		return nil, err
	}
	defer file.Close()
	return goquery.NewDocumentFromReader(file)
}

func (s *ScrapingServiceUkishimaTibaImpl) fetchTemperature(query string, doc *goquery.Document, ocean *model.Ocean) error {
	articleHTML, _ := doc.Find(query).Html()
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

func (s *ScrapingServiceUkishimaTibaImpl) fetchVisibility(query string, doc *goquery.Document, ocean *model.Ocean) error {
	articleHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`>透明度[\s\S]*ｍ</`)
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

func (s *ScrapingServiceUkishimaTibaImpl) fetchMeasuredTime(query string, doc *goquery.Document, ocean *model.Ocean) error {
	HTML := doc.Find(query)
	date, _ := HTML.Attr("datetime")
	reg := regexp.MustCompile(`[0-9０-９]{1,4}`)
	dates := reg.FindAllStringSubmatch(date, -1)

	// HACK
	year, err := strconv.Atoi(dates[0][0])
	month, err := strconv.Atoi(dates[1][0])
	day, err := strconv.Atoi(dates[2][0])
	hour, err := strconv.Atoi(dates[3][0])
	min, err := strconv.Atoi(dates[4][0])
	sec, err := strconv.Atoi(dates[5][0])
	if err != nil {
		s.ScrapingService.logging.Info("datetimeの変換に失敗")
		return err
	}
	ocean.MeasuredTime = time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
	return nil
}
