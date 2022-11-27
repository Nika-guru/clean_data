package crawl

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/log"
	dto "review-service/service/review/model/dto/revain"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductReviewsByPage(endpointDetail dto.EndpointDetail) error {
	for pageIdx := 1; ; pageIdx++ {
		url := getProductReviewsUrlFromEndpoint(endpointDetail.Endpoint, pageIdx)

		dom, err := GetHtmlDomByUrl(url)
		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_info.go/CrawlProductReviews/GetHtmlDomByUrl`, err.Error())
		}

		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			break
		}

		// ############# crawl detail comment #############
		productReviewRepo := extractProductReviewsByHtmlDom(dom, endpointDetail)
		CrawlProductReviewsInCurrentPage(productReviewRepo)
	}
	return nil
}

// all detail review of one product
func CrawlProductReviewsInCurrentPage(productReviewRepo dto.ProductReviewRepo) {
	for _, productReview := range productReviewRepo.ProductReviews {
		url := getProductReviewUrlFromEndpoint(productReview.Endpoint)
		dom, err := GetHtmlDomByUrl(url)
		if err != nil {
			log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_info.go/CrawlProductReview/GetHtmlDomByUrl`, err.Error())
		}

		//reponse not equal 200(404 --> No data to crawl)
		if dom == nil {
			continue //next detail review
		}

		//######## current detail project #########
		productReviewTmp := productReview
		extractProductReviewByHtmlDom(dom, productReviewTmp)
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

	domKey := `main` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 khpoB`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		//header(review, name, image ...)
		domKey := `div` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 bnwXZr`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey := `img` + ConvertClassesFormatFromBrowserToGoQuery(`LazyImage-sc-synjzy-0 ReviewTargetLogo__Logo-sc-160quaj-0 HAIBu gWKbJj`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				// val, ok := s.Attr(`data-src`)
				// fmt.Println("imageurl===", val, ok)
			})

			domKey = `div` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 gmjrOf`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				// fmt.Println("=====rating", s.Text())
			})

		})

		//footer(review)
		domKey = `article` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 Review__ReviewCard-sc-1xpzhiw-0 iQbXxL cDlhpG`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `a` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 kKJEOJ dDxbNj`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				endpointReview, foundEndpoint := s.Attr(`href`)
				if foundEndpoint {
					productReviewRepo.ProductReviews = append(productReviewRepo.ProductReviews, &dto.ProductReview{
						Endpoint:  endpointReview,
						ProductId: endpointDetail.ProductId,
					})
				}
			})

		})

	})

	return productReviewRepo
}

func extractProductReviewByHtmlDom(dom *goquery.Document, productReview *dto.ProductReview) {
	domKey := `script`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr(`type`)
		if ok && val == `application/ld+json` {

			var data any
			json.Unmarshal([]byte(s.Text()), &data)

			// headline, found := data.(map[string]any)[`headline`]
			// if found {
			// 	fmt.Println("headline", headline)
			// }

			// reviewBody, found := data.(map[string]any)[`reviewBody`]
			// if found {
			// 	fmt.Println("reviewBody", reviewBody)
			// }

			positiveNotes, found := data.(map[string]any)[`positiveNotes`]
			if found {
				fmt.Println(`======itemListElement`, positiveNotes.(map[string]any)[`itemListElement`].([]any)[0].(map[string]any)[`name`])
			}

			negativeNotes, found := data.(map[string]any)[`negativeNotes`]
			if found {
				fmt.Println(`======itemListElement`, negativeNotes.(map[string]any)[`itemListElement`].([]any)[0].(map[string]any)[`name`])
			}

			// author, found := data.(map[string]any)[`author`]
			// if found {
			// 	fmt.Println(`author name`, author.([]any)[0].(map[string]any)[`name`])
			// 	fmt.Println(`author name`, author.([]any)[0].(map[string]any)[`image`])
			// }

		}
	})

	//Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 ReviewPage__ReviewCard-sc-afh5kv-1 eUPSVm cTilou
	//	//Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 kMpeow  ---> name
	//	//Text-sc-kh4piv-0 jfONyJ ---> place dubai
}
