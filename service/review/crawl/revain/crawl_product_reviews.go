package crawl

import (
	"fmt"
	"review-service/pkg/log"
	dto "review-service/service/review/model/dto/revain"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductReviews(endpointDetail dto.EndpointDetail) error {
	for pageIdx := 1; pageIdx < 2; pageIdx++ {
		url := getProductReviewsUrlFromEndpoint(endpointDetail.Endpoint, pageIdx)

		dom, err := getHtmlDomByUrl(url)
		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_info.go/CrawlProductsInfoByProductType/getHtmlDomByUrl`, err.Error())
		}

		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			break
		}

		// ############# crawl detail comment #############
		productReviewRepo := extractProductReviewsByHtmlDom(dom, endpointDetail)
		CrawlProductReview(productReviewRepo)
	}
	return nil
}

func CrawlProductReview(productReviewRepo dto.ProductReviewRepo) {
	for _, productReview := range productReviewRepo.ProductReviews {
		url := getProductReviewUrlFromEndpoint(productReview.Endpoint)
		dom, err := getHtmlDomByUrl(url)
		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_info.go/CrawlProductReview/getHtmlDomByUrl`, err.Error())
		}

		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			continue //next detail review
		}
	}
}

func getProductReviewsUrlFromEndpoint(endpointProductInfo string, pageIdx int) string {
	params := fmt.Sprintf(_urlQueryParamsProductReviews, pageIdx) //bind data value to url param(s)
	url := (_baseUrl + endpointProductInfo + params)
	return url
}

func getProductReviewUrlFromEndpoint(endpointProductInfo string) string {
	url := (_baseUrl + endpointProductInfo)
	return url
}

func extractProductReviewsByHtmlDom(dom *goquery.Document, endpointDetail dto.EndpointDetail) dto.ProductReviewRepo {
	productReviewRepo := dto.ProductReviewRepo{}

	domKey := `main` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 khpoB`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		//header(review, name, image ...)
		domKey := `div` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 bnwXZr`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey := `img` + convertClassesFormatFromBrowserToGoQuery(`LazyImage-sc-synjzy-0 ReviewTargetLogo__Logo-sc-160quaj-0 HAIBu gWKbJj`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				// val, ok := s.Attr(`data-src`)
				// fmt.Println("imageurl===", val, ok)
			})

			domKey = `div` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 gmjrOf`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				// fmt.Println("=====rating", s.Text())
			})

		})

		//footer(review)
		domKey = `article` + convertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 Review__ReviewCard-sc-1xpzhiw-0 iQbXxL cDlhpG`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `a` + convertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 kKJEOJ dDxbNj`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				endpointReview, foundEndpoint := s.Attr(`href`)
				if foundEndpoint {
					productReviewRepo.ProductReviews = append(productReviewRepo.ProductReviews, dto.ProductReview{
						Endpoint:  endpointReview,
						ProductId: endpointDetail.ProductId,
					})
				}
			})

		})

	})

	return productReviewRepo
}

func extractProductReviewByHtmlDom(dom *goquery.Document) error {
	return nil
}
