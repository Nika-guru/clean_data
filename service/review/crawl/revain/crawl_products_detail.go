package crawl

import (
	"encoding/json"
	"fmt"
	"review-service/service/review/model/dao"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductDetail() {
	repo := dao.ProductRepo{}
	//TODO: get cache here
	//convert list to map
	productInfoMap := make(map[int64]dao.Product)
	for _, productInfo := range repo.Products {
		tmpProductInfo := *productInfo
		_, found := productInfoMap[int64(productInfo.ProductId)]
		if !found {
			productInfoMap[int64(productInfo.ProductId)] = tmpProductInfo
		}
	}
	for _, productInfo := range repo.Products {
		tmpProductInfo := *productInfo
		CrawlProductsDetailByProductInfo(tmpProductInfo, productInfoMap)
	}
}

func CrawlProductsDetailByProductInfo(productInfo dao.Product, productInfoMap map[int64]dao.Product) {
	// url := (_baseUrl + productInfo.EndpointProductDetail)

	// dom, err := getHtmlDomByUrl(url)
	// if err != nil {
	// 	log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_detail.go/CrawlProductsDetailByEndpoint/getHtmlDomByUrl`, err.Error())
	// }

	// // convert list to map

	// data := extractProductsDetailByHtmlDom(dom)
	// repo := dao.ProductDetailRepo{}
	// repo.ProductDetailList = data
	// fmt.Println(len(data))
	// repo.InsertDB(productInfoMap)
}

func extractProductsDetailByHtmlDom(dom *goquery.Document) []dao.ProductDetail {
	productDetailList := make([]dao.ProductDetail, 0)

	domKey := `main` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 khpoB`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		//header(review, name, image ...)
		domKey := `div` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 bnwXZr`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey := `img` + convertClassesFormatFromBrowserToGoQuery(`LazyImage-sc-synjzy-0 ReviewTargetLogo__Logo-sc-160quaj-0 HAIBu gWKbJj`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				val, ok := s.Attr(`data-src`)
				fmt.Println(val, ok)

			})

			domKey = `div` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 gmjrOf`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				fmt.Println(s.Text())
			})

		})

		//footer(review)
		domKey = `article` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 Review__ReviewCard-sc-1xpzhiw-0 iQbXxL cDlhpG`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `a` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 kKJEOJ dDxbNj`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				attrVal, foundAttr := s.Attr(`href`)
				if foundAttr {
					fmt.Println(`===review`, attrVal)
				}
			})

		})

	})

	//Body (get description, Official website, Social media )
	domKey = `script`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr(`type`)
		if ok && val == `application/ld+json` {

			var data any
			json.Unmarshal([]byte(s.Text()), &data)
			description, found := data.(map[string]any)["description"]
			if found {
				fmt.Println(`===description`, description)
			}

			url, found := data.(map[string]any)["url"]
			if found {
				fmt.Println(`===Official website`, url)
			}

			sameAs, found := data.(map[string]any)["sameAs"]
			if found {
				fmt.Println(`===Social media`, sameAs)
			}
		}
	})

	return productDetailList
}
