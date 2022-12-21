package utils

import (
	"crawler/pkg/log"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

const (
	WAIT_DURATION_WHEN_RATE_LIMIT = 5 * time.Second
	RESP_SUCCESS_STATUS_CODE      = 200
	RESP_NOT_FOUND_STATUS_CODE    = 404
)

func GetHtmlDomByUrl(url string) *goquery.Document {
beginCallAPI:
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Println(log.LogLevelDebug, `pkg/utils/crawl_html.go/GetHtmlDomByUrl/GetHtmlDomByUrl/http.Get(url)`, err.Error())
		time.Sleep(WAIT_DURATION_WHEN_RATE_LIMIT)
		goto beginCallAPI
	}
	defer res.Body.Close()

	if res.StatusCode != RESP_SUCCESS_STATUS_CODE {

		if RESP_NOT_FOUND_STATUS_CODE != res.StatusCode {
			time.Sleep(WAIT_DURATION_WHEN_RATE_LIMIT)
			goto beginCallAPI
		}

		//TODO: research why yobit got exception
		log.Println(log.LogLevelDebug, `pkg/utils/crawl_html.go/GetHtmlDomByUrl/GetHtmlDomByUrl/ RESP_NOT_FOUND_STATUS_CODE`, errors.New(`status respose not 200, actually: `+res.Status+` at url `+url))
		return nil
	}
	// Load the HTML document
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(log.LogLevelDebug, `pkg/utils/crawl_html.go/GetHtmlDomByUrl/GetHtmlDomByUrl/ goquery.NewDocumentFromReader(res.Body)`, err.Error())
		time.Sleep(WAIT_DURATION_WHEN_RATE_LIMIT)
		goto beginCallAPI
	}

	return dom
}

func GetHtmlDomJsRenderByUrl(url string) *goquery.Document {
	var dom *goquery.Document

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			dom = r.HTMLDoc
		},
		//BrowserEndpoint: "ws://localhost:3000",
	}).Start()

	return dom
}

func ConvertClassesFormatFromBrowserToGoQuery(input string) string {
	classes := input
	classes = `.` + classes
	classes = strings.ReplaceAll(classes, ` `, `.`)
	return classes
}
