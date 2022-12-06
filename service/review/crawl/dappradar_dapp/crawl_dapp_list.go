package crawl_dappradar_dapp

import (
	"review-service/pkg/utils"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlDappEndpointByHtmlDom(dom *goquery.Document) []dto_dappradar.EndpointDapp {
	dappradarList := make([]dto_dappradar.EndpointDapp, 0)

	domKey := `table` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-ehvNnt byUOet`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey := `tbody` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-laZRCg cblYxM`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `tr` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-eJDSGI gKTXmp`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				// fmt.Println(`run`)

				isAd := false
				domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-gfbthv egscaG`)
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {

					_txtAdDisplayHtml := `Ad`
					if s.Text() == _txtAdDisplayHtml {
						isAd = true
					}

				})

				if !isAd {
					dtoEndpointDapp := dto_dappradar.EndpointDapp{}
					dtoDetailDapp := dto_dappradar.DetailDapp{}
					dtoEndpointDapp.DetailDapp = &dtoDetailDapp

					//Image product
					domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-iAbQMe dBiXhX`)
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {

						domKey := `img`
						s.Find(domKey).Each(func(i int, s *goquery.Selection) {
							imgUrl, foundImgUrl := s.Attr(`src`)
							if foundImgUrl {
								dtoDetailDapp.Image = imgUrl
							}
						})

					})

					//Id, name product
					domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-hrDvXV bRkjlF`)
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {

						domKey := `a`
						s.Find(domKey).Each(func(i int, s *goquery.Selection) {
							endpointDetailUrl, foundEndpointDetailUrl := s.Attr(`href`)
							if foundEndpointDetailUrl {
								dtoEndpointDapp.Endpoint = endpointDetailUrl

								urlParts := strings.Split(endpointDetailUrl, `/`)

								if len(urlParts) > 3 {
									productId := urlParts[3]
									dtoDetailDapp.ProductId = productId

									productCategory := urlParts[2]
									dtoDetailDapp.CategoryId = productCategory

									productBlockchainId := urlParts[1]
									dtoDetailDapp.BlockchainId = productBlockchainId
								}

							}

							dtoDetailDapp.ProductName = s.Text()
						})

					})

					dappradarList = append(dappradarList, dtoEndpointDapp)
				}

			})

		})

	})
	return dappradarList
}
