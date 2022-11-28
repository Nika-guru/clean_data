package utils

import (
	"errors"
	"net/http"
	"review-service/service/constant"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetHtmlDomByUrl(url string) (*goquery.Document, error) {
beginCallAPI:
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != constant.RESP_SUCCESS_STATUS_CODE {

		if constant.RESP_TOO_MANY_REQ_STATUS_CODE == res.StatusCode {
			time.Sleep(constant.WAIT_DURATION_WHEN_RATE_LIMIT)
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
