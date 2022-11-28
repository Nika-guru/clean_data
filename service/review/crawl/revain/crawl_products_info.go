package crawl

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	dto "review-service/service/review/model/dto/revain"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductsInfo() {
	for _, endpointProductInfo := range constant.ENDPOINTS_PRODUCT_INFO_REVAIN {
		CrawlProductsInfoByProductType(endpointProductInfo)
	}

	log.Println(log.LogLevelInfo, `Crawl all from revain done`, `Crawl all from revain done`)
}

func CrawlProductsInfoByProductType(endpointProductInfo string) {
	for pageIdx := 1; pageIdx < 2; pageIdx++ {
		url := getProductsInfoUrlFromEndpoint(endpointProductInfo, pageIdx)

		dom, err := GetHtmlDomByUrl(url)
		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_info.go/CrawlProductsInfoByProductType/GetHtmlDomByUrl`, err.Error())
		}

		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			break
		}

		productInfoList := extractProductsInfoByHtmlDom(dom, url)
		productDtoRepo := &dto.ProductInfoRepo{}
		productDtoRepo.Products = productInfoList
		productDaoRepo := dao.ProductRepo{}
		productDaoRepo.ConverFrom(productDtoRepo)
		productDaoRepo.InsertDB() //insert finish --> has id incremental in model dao also dto

		productCategoryDaoRepo := dao.ProductCategoryRepo{}
		productCategoryDaoRepo.ConverFrom(productDtoRepo)
		productCategoryDaoRepo.InsertDB()

		//############# Crawl detail ###################
		endpointDetailRepo := dto.EndpointDetailRepo{}
		endpointDetailRepo.ConvertFrom(productDtoRepo)

		maxGoroutines := 10
		guard := make(chan struct{}, maxGoroutines)

		//SAve to cache, for crawl detail later
		for index, detailEndpoint := range endpointDetailRepo.Endpoints {
			guard <- struct{}{} // would block if guard channel is already filled
			go func(detailEndpoint dto.EndpointDetail, index int) {
				fmt.Println(`start index`, index)

				//Call detail product
				err := CrawlProductDetail(detailEndpoint)
				if err != nil {
					log.Println(log.LogLevelDebug, "service/review/crawl/revain/crawl_products_info/go/CrawlProductsInfoByProductType/CrawlProdcutDetail", err.Error())
					// continue
				}

				err = CrawlProductReviewsByPage(detailEndpoint)
				if err != nil {
					log.Println(log.LogLevelDebug, "service/review/crawl/revain/crawl_products_info/go/CrawlProductsInfoByProductType/CrawlProdcutReviews", err.Error())
					// continue
				}

				fmt.Println(`end index`, index)
				<-guard
			}(detailEndpoint, index)

		}

	}
}

func getProductsInfoUrlFromEndpoint(endpointProductInfo string, pageIdx int) string {
	params := fmt.Sprintf(_urlQueryParamsProductsInfo, pageIdx) //bind data value to url param(s)
	url := (_baseUrl + endpointProductInfo + params)
	return url
}

func extractProductsInfoByHtmlDom(dom *goquery.Document, currentUrl string) []*dto.ProductInfo {
	products := make([]*dto.ProductInfo, 0)

	domKey := `div` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 dzjLTP`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 ReviewTargetCard__Card-sc-qbvmhm-0 jMDOvK kHppZh`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			product := &dto.ProductInfo{}
			product.CrawlSource = currentUrl

			domKey = `img`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `data-src`
				imageUrl, foundAttrVal := s.Attr(attrKey)
				if foundAttrVal {
					product.ProductImage = imageUrl
				}

			})

			domKey = `a` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 gtFTPK gOcOhU`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `href`
				productDetailEndpoint, foundAttrVal := s.Attr(attrKey)
				if foundAttrVal {
					product.EndpointProductDetail = productDetailEndpoint
				}

				title := s.Text()
				product.ProductName = title

			})

			domKey = `div` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 bkUBKu`)
			typeNames := ``
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				typeNames += s.Text()
				if i < s.Length()-1 {
					typeNames += `, `
				}

			})
			typeNameArr := strings.Split(strings.Trim(typeNames, " "), `,`)
			product.ProductCategories = typeNameArr

			domKey = `p` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 ReviewTargetCard__LineClamp-sc-qbvmhm-1 ReviewTargetCard___StyledLineClamp-sc-qbvmhm-2 jVbmuR dUpQvL dxYXvO`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				shortDescription := s.Text()
				productDetail := make(map[string]any, 0)
				productDetail["shortDescription"] = shortDescription
				product.ProductDetail = productDetail

			})

			products = append(products, product)
		})

	})

	return products
}
