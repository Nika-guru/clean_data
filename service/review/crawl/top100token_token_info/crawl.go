package crawl_top100token_token_info

import (
	"fmt"
	"review-service/pkg/utils"
	"review-service/service/constant"

	"github.com/PuerkitoBio/goquery"
)

func Crawl() {
	url := fmt.Sprintf(`%s%s`, constant.BASE_URL_TOP100TOKEN, constant.ENDPOINT_LATEST_TOP100TOKEN)

	fmt.Println(url)
	dom := utils.GetHtmlDomByUrl(url)

	GetListchainByHtmlDom(dom)
}

func GetListchainByHtmlDom(dom *goquery.Document) {
	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`highlights`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		if i == 0 {

			domKey := `div`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				domKey := `span`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {

					fmt.Println(`====run here`, s.Text())
				})
			})

		}

	})
}
