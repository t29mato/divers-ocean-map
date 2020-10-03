package iop

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

// FetchServiceImpl ...
type FetchServiceImpl struct {
	logging           *logging.OceanLoggingImpl
	name              string
	url               string
	queryTemperature  string
	queryVisibility   string
	queryMeasuredTime string
}

// NewFetchService ...
func NewFetchService(name string, url string, logging *logging.OceanLoggingImpl) *FetchServiceImpl {
	// HACK: FetchServiceImplの中で、各ダイビングポイントの場所に依存して変わるもの以外は、ダイビングポイントの構造体に直接持たせる (隠したい)
	return &FetchServiceImpl{
		logging:           logging,
		name:              name,
		url:               url,
		queryTemperature:  "#homeConditionDetail > dl > dd:nth-child(2)",
		queryVisibility:   "#homeConditionDetail > dl > dd:nth-child(4)",
		queryMeasuredTime: "#homeConditionIndex > dl > dt",
	}
}

// Fetch ...
func (s *FetchServiceImpl) Fetch() (*model.Ocean, error) {
	s.logging.Info("伊豆海洋公園のスクレイピング開始")
	ocean := model.NewOcean(s.name, s.url)

	// DOM取得
	doc, err := s.fetchDocument(s.url)
	if err != nil {
		s.logging.Info("HTMLファイルの読み込みに失敗しました。url =", s.url)
		return nil, err
	}

	// 水温取得
	err = s.fetchTemperature(s.queryTemperature, doc, ocean)
	if err != nil {
		s.logging.Info("水温の取得に失敗")
		return nil, err
	}

	// 透明度取得
	err = s.fetchVisibility(s.queryVisibility, doc, ocean)
	if err != nil {
		s.logging.Info("透明度の取得に失敗")
		return nil, err
	}

	// 測定日時取得
	err = s.fetchMeasuredTime(s.queryMeasuredTime, doc, ocean)
	if err != nil {
		s.logging.Info("測定日時の取得に失敗")
		return nil, err
	}

	s.logging.Info("伊豆海洋公園のスクレイピング終了")
	return ocean, err
}

func (s *FetchServiceImpl) fetchDocument(url string, ocean *model.Ocean) (*goquery.Document, error) {
	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if strings.Contains(url, "http") {
		return goquery.NewDocument(url)
	}
	file, err := os.Open(url)
	if err != nil {
		s.logging.Info("有効なファイルパスでありません。url =", url)
		return nil, err
	}
	defer file.Close()
	return goquery.NewDocumentFromReader(file)
}

func (s *FetchServiceImpl) fetchTemperature(query string, doc *goquery.Document, ocean *model.Ocean) error {
	temperatureHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	temperatures := reg.FindAllStringSubmatch(temperatureHTML, -1)
	min, err := strconv.Atoi(temperatures[0][0])
	if err != nil {
		s.logging.Info("水温1の数値変換に失敗, 変換前=", temperatures[0][0])
		return err
	}
	switch len(temperatures) {
	case 1:
		ocean.Temperature.Med = min
	case 2:
		max, err := strconv.Atoi(temperatures[1][0])
		if err != nil {
			s.logging.Info("水温2の数値変換に失敗, 変換前=", temperatures[1][0])
			return err
		}
		ocean.Temperature.Min = min
		ocean.Temperature.Max = max
	}
	return nil
}

func (s *FetchServiceImpl) fetchVisibility(query string, doc *goquery.Document, ocean *model.Ocean) error {
	visibilityHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	visibilities := reg.FindAllStringSubmatch(visibilityHTML, -1)

	min, err := strconv.Atoi(visibilities[0][0])
	if err != nil {
		s.logging.Info("透明度1の数値変換に失敗, 変換前=", visibilities[0][0])
		return err
	}
	switch len(visibilities) {
	case 1:
		ocean.Visibility.Med = min
	case 2:
		max, err := strconv.Atoi(visibilities[1][0])
		if err != nil {
			s.logging.Info("透明度2の数値変換に失敗, 変換前=", visibilities[1][0])
			return err
		}
		ocean.Visibility.Min = min
		ocean.Visibility.Max = max
	}
	return nil
}

func (s *FetchServiceImpl) fetchMeasuredTime(query string, doc *goquery.Document, ocean *model.Ocean) error {
	MeasuredTimeHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`\d{1,2}`)
	measuredTimes := reg.FindAllStringSubmatch(MeasuredTimeHTML, -1)

	month, err := strconv.Atoi(measuredTimes[0][0])
	if err != nil {
		s.logging.Info("月の数値変換に失敗, 変換前=", measuredTimes[0][0])
		return err
	}
	day, err := strconv.Atoi(measuredTimes[1][0])
	if err != nil {
		s.logging.Info("日の数値変換に失敗, 変換前=", measuredTimes[1][0])
		return err
	}
	ocean.MeasuredTime = time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return nil
}
