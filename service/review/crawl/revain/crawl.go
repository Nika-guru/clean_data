package crawl

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	_baseUrl        = `https://revain.org`
	_urlQueryParams = `?sortBy=reviews&page=%d`
)

func init() {
	CrawlProductsInfo()
	CrawlProductDetail()
}
func getHtmlDomByUrl(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
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

func convertClassesFormatFromBrowserToGoQuery(input string) string {
	classes := input
	classes = `.` + classes
	classes = strings.ReplaceAll(classes, ` `, `.`)
	return classes
}
