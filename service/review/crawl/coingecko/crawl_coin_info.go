package crawl_coingecko

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto_coingecko "review-service/service/review/model/dto/coingecko"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductIdByCategory(endpointCategory *dto_coingecko.EndpointCategory) {
	for pageIdx := 1; pageIdx < 2; pageIdx++ {
		url := getUrlProductIdByCategory(endpointCategory.Endpoint, pageIdx)
		dom, err := utils.GetHtmlDomByUrl(url)

		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/coingecko/crawl_coin_info.go/CrawlProductNameById/GetHtmlDomByUrl`, err.Error())
		}

		//TODO: check by another way
		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			return
		}

		data := extractProductIdByHtmlDom(dom)
		endpointCategory.CoinIdList = data
	}
}

func getUrlProductIdByCategory(endpoint string, pageIdx int) string {
	params := fmt.Sprintf(_paramsProductIdByCategory, pageIdx)
	url := (_baseUrl + endpoint + params)
	return url
}

func extractProductIdByHtmlDom(dom *goquery.Document) []string {
	coinIdList := make([]string, 0)

	domKey := `table` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sort table mb-0 text-sm text-lg-normal table-scrollable`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `tbody`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `tr`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`tw-flex-auto`)
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {

					domKey = `a`
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {

						attrKey := `href`
						urlDetail, foundUrlDetail := s.Attr(attrKey)
						if foundUrlDetail {
							urlParts := strings.Split(urlDetail, `/`)
							coinId := urlParts[len(urlParts)-1]

							coinIdList = append(coinIdList, coinId)
						}

					})

				})

			})

		})

	})

	return coinIdList
}
