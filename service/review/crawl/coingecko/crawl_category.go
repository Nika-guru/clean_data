package crawl_coingecko

import (
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	dto_coingecko "review-service/service/review/model/dto/coingecko"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductCategories() {
	dom := utils.GetHtmlDomByUrl(`https://www.coingecko.com/en/categories`)

	//reponse not equal 200(404 --> No data to crawl)
	if dom == nil {
		return
	}

	endpointCategories := extractCategoryByHtmlDom(dom)

	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)
	endpointCategories = endpointCategories[0:2]
	routineCount := len(endpointCategories)
	done := make(chan int, routineCount)

	for index, endpointCategory := range endpointCategories {
		guard <- struct{}{} //buffered channel, full capacity, wait here --> limit go routine

		go func(endpointCategory *dto_coingecko.EndpointCategory, index int) {

			CrawlProductIdByCategory(endpointCategory) //after run here, data updated, list coinId

			<-guard
			done <- index
		}(endpointCategory, index)
	}

	//Wait all go routine done
	for i := 0; i < len(endpointCategories); i++ {
		index := <-done

		//Data here
		// fmt.Println(endpointCategories[index].CategoryName, endpointCategories[index].CoinIdList)
		// fmt.Println(index, i, `================================================================`)

		//################ Insert not duplicated default category #################
		category := dao.Category{
			CategoryName: constant.DEFAULT_CATEGORY_PRODUCT_REVAIN,
		}
		isExist, err := category.SelectByName()
		if err != nil {
			log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/category.SelectByName()`, err.Error())
			return
		}

		if !isExist {
			err = category.InsertDB()
			if err != nil {
				log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/category.InsertDB()`, err.Error())
				return
			}
		}

		endpointCategory := endpointCategories[index]

		//################ Insert not duplicated sub category #################
		subcategory := dao.SubCategory{
			CategoryId:      category.CategoryId,
			SubCategoryName: endpointCategory.CategoryName,
		}
		isExist, err = subcategory.SelectByName()
		if err != nil {
			log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/subcategory.SelectByName()`, err.Error())
			return
		}
		if !isExist {
			err = subcategory.InsertDB()
			if err != nil {
				log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/subcategory.InsertDB()`, err.Error())
				return
			}
		}

		for _, coinId := range endpointCategory.CoinIdList {
			//################ Insert not duplicated product name #################
			product := dao.Product{
				ProductName: coinId,
			}
			isExist, err = product.SelectByProductName()
			if err != nil {
				log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/product.SelectByProductName()`, err.Error())
				continue
			}
			if !isExist {
				err = product.InsertDB()
				if err != nil {
					log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/product.InsertDB()`, err.Error())
					continue
				}
			}

			productCategory := dao.ProductCategory{
				CategoryId:    category.CategoryId,
				SubCategoryId: &subcategory.SubCategoryId,
				ProductId:     product.ProductId,
			}
			err = productCategory.InsertDB()
			if err != nil {
				log.Println(log.LogLevelError, `service/review/crawl/coingecko/crawl_category.go/CrawlProductCategories/productCategory.InsertDB()`, err.Error())
				continue
			} else {
				// fmt.Println(`run here`)
			}

		}

	}

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
