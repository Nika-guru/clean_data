package crawl_revain

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	dto_revain "review-service/service/review/model/dto/revain"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductsInfo() {
	log.Println(log.LogLevelInfo, `Start crawl revain`, `Start crawl revain`)

	for _, endpointProductInfo := range constant.ENDPOINTS_PRODUCT_INFO_REVAIN {
		log.Println(log.LogLevelInfo, `Start crawl revain/ endpoint `+endpointProductInfo, `Start crawl revain/ endpoint `+endpointProductInfo)

		CrawlProductsInfoByProductType(endpointProductInfo)

		log.Println(log.LogLevelInfo, `End crawl revain/ endpoint `+endpointProductInfo, `End crawl revain/ endpoint `+endpointProductInfo)
	}

	log.Println(log.LogLevelInfo, `End crawl revain`, `End crawl revain`)
}

func CrawlProductsInfoByProductType(endpointProductInfo string) {
	for pageIdx := 1; ; pageIdx++ {
		url := getProductsInfoUrlFromEndpoint(endpointProductInfo, pageIdx)

		dom := utils.GetHtmlDomByUrl(url)

		//reponse not equal 200, 404(truy to call) --> No data to crawl
		if dom == nil {
			break
		}

		productInfoList := extractProductsInfoByHtmlDom(dom, url)
		productDtoRepo := &dto_revain.ProductInfoRepo{}
		productDtoRepo.Products = productInfoList

		productDaoRepo := dao.ProductRepo{}
		productDaoRepo.ConverFrom(productDtoRepo)
		//insert finish --> has id incremental in model dao also dto
		productDaoRepo.InsertDB(&dto_revain.ProductInfoDebug{
			EndpointProduct: endpointProductInfo,
			Url:             url,
			PageIndex:       uint8(pageIdx),
			IsSuccess:       false,
		})

		CrawlMoreDetailData(productDtoRepo)
	}
}

func CrawlMoreDetailData(productDtoRepo *dto_revain.ProductInfoRepo) {
	//############# Crawl detail and all its review ###################
	productCategoryDaoRepo := dao.ProductCategoryRepo{}
	productCategoryDaoRepo.ConverFrom(productDtoRepo)
	productCategoryDaoRepo.InsertDB()

	//############# Crawl detail and all its review ###################
	endpointDetailRepo := dto_revain.EndpointDetailRepo{}
	endpointDetailRepo.ConvertFrom(productDtoRepo)

	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)

	//Crawl detail and review
	for index, detailEndpoint := range endpointDetailRepo.Endpoints {
		guard <- struct{}{} // would block if guard channel is already filled
		go func(detailEndpoint dto_revain.EndpointDetail, index int) {
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
			<-guard
		}(detailEndpoint, index)

	}
}

func getProductsInfoUrlFromEndpoint(endpointProductInfo string, pageIdx int) string {
	params := fmt.Sprintf(_urlQueryParamsProductsInfo, pageIdx) //bind data value to url param(s)
	url := (_baseUrl + endpointProductInfo + params)
	return url
}

func extractProductsInfoByHtmlDom(dom *goquery.Document, currentUrl string) []*dto_revain.ProductInfo {
	products := make([]*dto_revain.ProductInfo, 0)

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 dzjLTP`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 ReviewTargetCard__Card-sc-qbvmhm-0 jMDOvK kHppZh`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			product := &dto_revain.ProductInfo{}
			product.CrawlSource = currentUrl

			domKey = `img`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `data-src`
				imageUrl, foundAttrVal := s.Attr(attrKey)
				if foundAttrVal {
					product.ProductImage = imageUrl
				}

			})

			domKey = `a` + utils.ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 gtFTPK gOcOhU`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `href`
				productDetailEndpoint, foundAttrVal := s.Attr(attrKey)
				if foundAttrVal {
					product.EndpointProductDetail = productDetailEndpoint

					detailEndpointParts := strings.Split(productDetailEndpoint, `/`)
					title := detailEndpointParts[len(detailEndpointParts)-1]
					product.ProductName = title
				}

				// productName := s.Text()

			})

			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 bkUBKu`)
			typeNames := ``
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				typeNames += s.Text()
				if i < s.Length()-1 {
					typeNames += `, `
				}

			})
			typeNameArr := strings.Split(strings.Trim(typeNames, " "), `,`)
			product.ProductCategories = typeNameArr

			domKey = `p` + utils.ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 ReviewTargetCard__LineClamp-sc-qbvmhm-1 ReviewTargetCard___StyledLineClamp-sc-qbvmhm-2 jVbmuR dUpQvL dxYXvO`)
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
