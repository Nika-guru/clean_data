package crawl

import (
	"encoding/json"
	"review-service/pkg/log"
	"review-service/service/review/model/dao"
	dto "review-service/service/review/model/dto/revain"

	"github.com/PuerkitoBio/goquery"
)

func CrawlProductDetail(endpoint dto.EndpointDetail) error {
	url := getProductDetailUrlFromEndpoint(endpoint.Endpoint)

	dom, err := GetHtmlDomByUrl(url)
	if err != nil {
		log.Println(log.LogLevelError, `review/crawl/revain/crawl_products_detail.go/CrawlProdcutDetail/GetHtmlDomByUrl`, err.Error())
	}

	//reponse not equal 200(404 --> No data to crawl)
	if dom == nil {
		return nil
	}

	dtoProductDetail := extractProductsDetailByHtmlDom(dom)
	dtoProductDetail.ProductId = endpoint.ProductId
	daoProduct := dao.Product{}
	daoProduct.ConvertFrom(dtoProductDetail)
	err = daoProduct.UpdateDescription()
	if err != nil {
		return err
	}

	productContactRepo := dao.ProductContactRepo{}
	productContactRepo.ConvertFrom(dtoProductDetail)
	productContactRepo.InsertDB()

	return nil
}

func getProductDetailUrlFromEndpoint(endpointProductDetail string) string {
	url := (_baseUrl + endpointProductDetail)
	return url
}

func extractProductsDetailByHtmlDom(dom *goquery.Document) dto.ProductDetail {
	productDetail := dto.ProductDetail{}
	if productDetail.Url == nil {
		productDetail.Url = make(map[string]string)
	}
	//Body (get description, Official website, Social media )
	domKey := `script`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr(`type`)
		if ok && val == `application/ld+json` {

			var data any
			json.Unmarshal([]byte(s.Text()), &data)
			description, found := data.(map[string]any)["description"]
			if found {
				productDetail.Description = description.(string)
			}

			url, found := data.(map[string]any)["url"]
			if found {
				productDetail.Url[url.(string)] = `Official website`
			}

			urlArr, found := data.(map[string]any)["sameAs"]
			if found {
				for _, url := range urlArr.([]any) {
					productDetail.Url[url.(string)] = `Social media`
				}
			}
		}
	})

	return productDetail
}
