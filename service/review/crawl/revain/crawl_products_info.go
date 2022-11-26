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
}

func CrawlProductsInfoByProductType(endpointProductInfo string) {
	for pageIdx := 1; pageIdx < 2; pageIdx++ {
		url := getProductsInfoUrlFromProductTypeName(endpointProductInfo, pageIdx)

		dom, err := getHtmlDomByUrl(url)
		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_info.go/CrawlProductsInfoByProductType/getHtmlDomByUrl`, err.Error())
		}

		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			break
		}

		productInfo := extractProductsInfoByHtmlDom(dom, url)
		productDtoRepo := &dto.ProductInfoRepo{}
		productDtoRepo.Products = productInfo
		productDaoRepo := dao.ProductRepo{}
		productDaoRepo.ConverFrom(productDtoRepo)
		productDaoRepo.InsertDB() //insert finish --> has id

		// endpointDetailRepo := dto.EndpointDetailRepo{}
		// //get list endpoint product id
		// for _, productInfo := range productRepo.Products {
		// 	endpointDetailRepo.Endpoints = append(endpointDetailRepo.Endpoints, dto.EndpointDetail{
		// 		ProductId: productInfo.ProductId,
		// 		Endpoint:  productInfo.EndpointProductDetail,
		// 	})

		// 	productInfo.ProductId
		// 	productInfo.ProductCategories
		// }
	}
}

func getProductsInfoUrlFromProductTypeName(endpointProductInfo string, pageIdx int) string {
	params := fmt.Sprintf(_urlQueryParams, pageIdx) //bind data value to url param(s)
	url := (_baseUrl + endpointProductInfo + params)
	return url
}

func extractProductsInfoByHtmlDom(dom *goquery.Document, currentUrl string) []dto.ProductInfo {
	products := make([]dto.ProductInfo, 0)

	domKey := `div` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 dzjLTP`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 ReviewTargetCard__Card-sc-qbvmhm-0 jMDOvK kHppZh`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			product := dto.ProductInfo{}
			product.CrawlSource = currentUrl

			domKey = `img`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `data-src`
				imageUrl, foundAttrVal := s.Attr(attrKey)
				if foundAttrVal {
					product.ProductImage = imageUrl
				}

			})

			domKey = `a` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 gtFTPK gOcOhU`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `href`
				productDetailEndpoint, foundAttrVal := s.Attr(attrKey)
				if foundAttrVal {
					product.EndpointProductDetail = productDetailEndpoint
				}

				title := s.Text()
				product.ProductName = title

			})

			domKey = `div` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 bkUBKu`)
			typeNames := ``
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				typeNames += s.Text()
				if i < s.Length()-1 {
					typeNames += `, `
				}

			})
			typeNameArr := strings.Split(strings.ReplaceAll(typeNames, " ", ""), `,`)
			product.ProductCategories = typeNameArr

			domKey = `p` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 ReviewTargetCard__LineClamp-sc-qbvmhm-1 ReviewTargetCard___StyledLineClamp-sc-qbvmhm-2 jVbmuR dUpQvL dxYXvO`)
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
