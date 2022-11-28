package crawl_revain

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/review/model/dao"
	dto "review-service/service/review/model/dto/revain"

	"github.com/PuerkitoBio/goquery"
)

// all detail review of one product
func CrawlProductReviewsInCurrentPage(productReviewRepo dto.ProductReviewRepo) {
	for _, productReview := range productReviewRepo.ProductReviews {
		url := getProductReviewUrlFromEndpoint(productReview.Endpoint)
		dom, err := utils.GetHtmlDomByUrl(url)
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

		account := dao.Account{}
		account.ConvertFrom(*productReviewTmp)
		err = account.InsertDB() //after insert have account id
		if err != nil {
			log.Println(log.LogLevelError, `service/review/crawl/revain/crawl_pagination_product_review.go/CrawlProductReviewsInCurrentPage/account.InsertDB()`, err.Error())
		}

		review := dao.Review{}
		review.ConvertFrom(*productReviewTmp)
		review.AccountId = account.Id
		err = review.InsertDB()
		if err != nil {
			log.Println(log.LogLevelError, `service/review/crawl/revain/crawl_pagination_product_review.go/CrawlProductReviewsInCurrentPage/review.InsertDB()`, err.Error())
		}

	}
}

func getProductReviewUrlFromEndpoint(endpointProductInfo string) string {
	url := (_baseUrl + endpointProductInfo)
	return url
}

func extractProductReviewByHtmlDom(dom *goquery.Document, productReview *dto.ProductReview) {
	domKey := `script`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr(`type`)
		if ok && val == `application/ld+json` {
			content := ``

			var data any
			json.Unmarshal([]byte(s.Text()), &data)

			headline, found := data.(map[string]any)[`headline`]
			if found {
				content += headline.(string)
			}

			reviewBody, found := data.(map[string]any)[`reviewBody`]
			if found {
				content += `\n`
				content += reviewBody.(string)
				content += `\n`
			}

			positiveNotes, found := data.(map[string]any)[`positiveNotes`]
			if found {
				content += `\n------------------\n`
				content += `- Positive notes: \n`
				for _, val := range positiveNotes.(map[string]any)[`itemListElement`].([]any) {
					content += fmt.Sprintf(`	+ %s\n`, val.(map[string]any)[`name`])
				}
			}

			negativeNotes, found := data.(map[string]any)[`negativeNotes`]
			if found {
				content += `\n------------------\n`
				content += `- Negative notes: \n`
				for _, val := range negativeNotes.(map[string]any)[`itemListElement`].([]any) {
					content += fmt.Sprintf(`	+ %s\n`, val.(map[string]any)[`name`])
				}
			}

			productReview.Content = content

			author, found := data.(map[string]any)[`author`]
			if found {
				username := author.([]any)[0].(map[string]any)[`name`].(string)
				productReview.Username = username

				image := author.([]any)[0].(map[string]any)[`image`].(string)
				productReview.AccountImage = image
			}

			reviewRating, found := data.(map[string]any)[`reviewRating`]
			if found {
				productReview.Star = reviewRating.(map[string]any)[`ratingValue`].(float64)
			}

			datePublished, found := data.(map[string]any)[`datePublished`]
			if found {
				time, err := utils.StringToTimestamp(datePublished.(string))
				if err != nil {
					log.Println(log.LogLevelFatal, `service/review/crawl/revain/crawl_pagination_product_review.go`, err.Error())
				}
				productReview.ReviewDate = utils.TimestampToString(time)
			}
		}
	})

	//Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 ReviewPage__ReviewCard-sc-afh5kv-1 eUPSVm cTilou
	//	//Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 kMpeow  ---> name
	//	//Text-sc-kh4piv-0 jfONyJ ---> place dubai

}
