package crawl_dappradar_dapp

import (
	"fmt"
	"review-service/pkg/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//TODO: check data length response equal zero
func CrawlDappDetailByHtmlDom(dom *goquery.Document) {
	//Tags
	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-hhNnGo kXsScc`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-hItEmJ jqLsUG`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey := `span`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				tagName := s.Text()
				fmt.Println(tagName)
			})

		})

	})

	//Social Media, (default type: Social media)
	// domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-ekFWYn ecxPlE`)
	// dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

	// 	domKey = `a`
	// 	s.Find(domKey).Each(func(i int, s *goquery.Selection) {
	// 		socialUrl, foundSocialUrl := s.Attr(`href`)
	// 		if foundSocialUrl {
	// 			socialImgSvg, err := s.Html()
	// 			if err == nil {
	// 				fmt.Println(socialUrl, socialImgSvg, `Social media`)
	// 			}
	// 		}
	// 	})

	// })

	//Descripttion
	domKey = `p`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		description := strings.TrimSpace(s.Text())
		if description != `` {
			fmt.Println(description)
		}
	})

}
