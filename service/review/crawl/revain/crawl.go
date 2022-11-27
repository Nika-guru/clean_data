package crawl

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	_baseUrl                      = `https://revain.org`
	_urlQueryParamsProductsInfo   = `?sortBy=reviews&page=%d`
	_urlQueryParamsProductReviews = `?page=%d&sortBy=recent&direction=ASC`
)

const (
	_respSuccessStatusCode     = 200
	_respTooManyReqStatusCode  = 429
	_waitDurationWhenRateLimit = 5 * time.Second
)

func init() {
	go CrawlProductsInfo()
}
func GetHtmlDomByUrl(url string) (*goquery.Document, error) {
beginCallAPI:
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != _respSuccessStatusCode {

		if _respTooManyReqStatusCode == res.StatusCode {
			time.Sleep(_waitDurationWhenRateLimit)
			goto beginCallAPI
		}

		//TODO: research why yobit got exception
		return nil, errors.New(`status respose not 200, actually: ` + res.Status + ` at url ` + url)
	}
	// Load the HTML document
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return dom, nil
}

func ConvertClassesFormatFromBrowserToGoQuery(input string) string {
	classes := input
	classes = `.` + classes
	classes = strings.ReplaceAll(classes, ` `, `.`)
	return classes
}
