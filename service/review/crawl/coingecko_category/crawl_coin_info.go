package crawl_coingecko_category

import (
	"fmt"
	"review-service/pkg/utils"
	dto_coingecko "review-service/service/review/model/dto/coingecko"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductIdByCategory(endpointCategory *dto_coingecko.EndpointCategory) {
	for pageIdx := 1; ; pageIdx++ {
		url := getUrlProductIdByCategory(endpointCategory.Endpoint, pageIdx)
		dom := utils.GetHtmlDomByUrl(url)

		//reponse not have data (table)
		if !IsExistCoinByDom(dom) {
			break
		}

		data := extractProductIdByHtmlDom(dom)
		endpointCategory.CoinIdList = append(endpointCategory.CoinIdList, data...)

		// #################################Start Debug #################################
		debug := dto_coingecko.Debug{}
		debug.AddProductCategory(dto_coingecko.ProductCategoryDebug{
			CategoryName: endpointCategory.CategoryName,

			Url:       url,
			PageIndex: uint8(pageIdx),
			IsSuccess: true,
		})
		// #################################End Debug #################################
	}
}

func IsExistCoinByDom(dom *goquery.Document) bool {
	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`coingecko-table`)

	isExist := false
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		isExist = true
	})
	return isExist
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
