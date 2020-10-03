package ukishima

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"scraping/logging"
	"scraping/model"
	"scraping/service/util"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// FetchServiceImpl ...
type FetchServiceImpl struct {
	logging      *logging.OceanLoggingImpl
	name         string
	url          string
	queryArticle string
	queryDate    string
}

// NewFetchService ...
// TODO: New functionは全て同じ関数名にする
func NewFetchService(name string, url string, logging *logging.OceanLoggingImpl) *FetchServiceImpl {
	// HACK: FetchServiceImplの中で、各ダイビングポイントの場所に依存して変わるもの以外は、ダイビングポイントの構造体に直接持たせる (隠したい)
	return &FetchServiceImpl{
		logging:      logging,
		name:         name,
		url:          url,
		queryArticle: "#content > div.blog > div > div.text",
		queryDate:    "#content > div:nth-child(2) > h2",
	}
}

// Fetch ...
func (s *FetchServiceImpl) Fetch() (*model.Ocean, error) {
	s.logging.Info("浮島(西伊豆)のスクレイピング開始")
	ocean := model.NewOcean(s.name, s.url)

	// DOM取得
	doc, err := s.fetchDocument(s.url, ocean)
	if err != nil {
		s.logging.Info("HTMLファイルの読み込みに失敗しました。url =", s.url)
		return nil, err
	}

	// 水温取得
	err = s.fetchTemperature(s.queryArticle, doc, ocean)
	if err != nil {
		s.logging.Info("水温の取得に失敗")
		return nil, err
	}

	// 透明度取得
	err = s.fetchVisibility(s.queryArticle, doc, ocean)
	if err != nil {
		s.logging.Info("透明度の取得に失敗")
		return nil, err
	}

	// 測定日時取得
	err = s.fetchMeasuredTime(s.queryDate, doc, ocean)
	if err != nil {
		s.logging.Info("測定日時の取得に失敗")
		return nil, err
	}

	s.logging.Info("浮島(西伊豆)のスクレイピング終了")
	return ocean, err
}

func (s *FetchServiceImpl) fetchDocument(url string, ocean *model.Ocean) (*goquery.Document, error) {
	// 単体テスト実行時はローカルのHTMLファイルから取得する
	if strings.Contains(url, "http") {
		doc, _ := goquery.NewDocument(url)
		// トップページには透明度情報がないので、トップページから最新の記事のURLを取得する
		latestArticleURL, _ := doc.Find("#content > div:nth-child(2) > div > h3 > a").Attr("href")
		ocean.URL = latestArticleURL

		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		utfBody := transform.NewReader(bufio.NewReader(res.Body), japanese.ShiftJIS.NewDecoder())
		latestDoc, err := goquery.NewDocumentFromReader(utfBody)
		if err != nil {
			return nil, err
		}

		return latestDoc, nil
	}
	file, err := os.Open(url)
	defer file.Close()

	reader := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())
	// 書き込み先ファイルを用意
	utf8File, err := os.Create("./utf-8.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer utf8File.Close()
	tee := io.TeeReader(reader, utf8File)
	scanner := bufio.NewScanner(tee)
	for scanner.Scan() {
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	file, err = os.Open("./utf-8.txt")
	defer os.Remove("./utf-8.txt")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return goquery.NewDocumentFromReader(file)
}

func (s *FetchServiceImpl) fetchTemperature(query string, doc *goquery.Document, ocean *model.Ocean) error {
	articleHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`水温[\s\S]*℃`)
	temperatureHTML := reg.FindAllStringSubmatch(articleHTML, -1)
	reg = regexp.MustCompile(`[0-9０-９]{1,2}`)

	if len(temperatureHTML) == 0 {
		s.logging.Info("水温を1箇所も取得に失敗 URL=", s.url)
		return nil
	}
	temperatures := reg.FindAllStringSubmatch(temperatureHTML[0][0], -1)
	fmt.Println(temperatureHTML)

	min, err := util.ConvertIntFromFullWidthString(&temperatures[0][0])
	if err != nil {
		s.logging.Info("水温1の数値変換に失敗, 変換前=", temperatures[0][0])
		return err
	}
	switch len(temperatures) {
	case 1:
		ocean.Temperature.Med = min
	case 2:
		max, err := util.ConvertIntFromFullWidthString(&temperatures[1][0])
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
	articleHTML, _ := doc.Find(query).Html()
	reg := regexp.MustCompile(`透明度[\s\S]*ｍ`)
	visibilityHTML := reg.FindAllStringSubmatch(articleHTML, -1)
	reg = regexp.MustCompile(`[0-9０-９]{1,2}`)

	if len(visibilityHTML) == 0 {
		s.logging.Info("透明度を1箇所も取得に失敗 URL=", s.url)
		return nil
	}
	visibilities := reg.FindAllStringSubmatch(visibilityHTML[0][0], -1)

	min, err := util.ConvertIntFromFullWidthString(&visibilities[0][0])
	if err != nil {
		s.logging.Info("透明度1の数値変換に失敗, 変換前=", visibilities[0][0])
		return err
	}
	switch len(visibilities) {
	case 1:
		ocean.Visibility.Med = min
	case 2:
		max, err := util.ConvertIntFromFullWidthString(&visibilities[1][0])
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
	reg := regexp.MustCompile(`\d{1,4}`)
	measuredTimes := reg.FindAllStringSubmatch(MeasuredTimeHTML, -1)

	// HACK
	year, err := strconv.Atoi(measuredTimes[0][0])
	month, err := strconv.Atoi(measuredTimes[1][0])
	day, err := strconv.Atoi(measuredTimes[2][0])
	if err != nil {
		s.logging.Info("datetimeの変換に失敗")
		return err
	}
	ocean.MeasuredTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return nil
}
