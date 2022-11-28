package crawl_coingecko

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto_coingecko "review-service/service/review/model/dto/coingecko"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductCategories() {
	dom, err := utils.GetHtmlDomByUrl(`https://www.coingecko.com/en/categories`)

	if err != nil {
		log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/GetHtmlDomByUrl`, err.Error())
	}

	//reponse not equal 200(404 --> No data to crawl)
	if dom == nil {
		return
	}

	endpointCategories := extractCategoryByHtmlDom(dom)

	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)
	routineCount := len(endpointCategories)
	done := make(chan struct{}, routineCount)

	for _, endpointCategory := range endpointCategories {
		guard <- struct{}{} //buffered channel, full capacity, wait here --> limit go routine

		go func(endpointCategory *dto_coingecko.EndpointCategory) {

			CrawlProductIdByCategory(endpointCategory) //after run here, data updated, list coinId

			<-guard
			done <- struct{}{}
		}(endpointCategory)
	}

	//Wait all go routine done
	for i := 0; i < len(endpointCategories); i++ {
		<-done
	}

	//Data here
	fmt.Println(endpointCategories[0].CategoryName, endpointCategories[0].CoinIdList)
}

func extractCategoryByHtmlDom(dom *goquery.Document) []*dto_coingecko.EndpointCategory {
	endpointCategories := make([]*dto_coingecko.EndpointCategory, 0)

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`gecko-table-container tw-mt-3`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `tbody`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `tr`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				endpointCategory := &dto_coingecko.EndpointCategory{}

				domKey = `b`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {

					domKey = `a`
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {
						categoryName := s.Text()
						endpointCategory.CategoryName = categoryName

						attrKey := `href`
						endpointListCoin, found := s.Attr(attrKey)
						if found {
							endpointCategory.Endpoint = endpointListCoin
						}

					})

				})

				endpointCategories = append(endpointCategories, endpointCategory)
			})

		})

	})

	return endpointCategories
}
