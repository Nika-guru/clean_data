package main

import (
	"fmt"
	crawl "review-service/service/review/crawl/revain"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	_baseUrl        = `https://revain.org`
	_productListUrl = (_baseUrl + `/exchanges?sortBy=reviews&page=%d`)
)

func main() {
	crawl.CrawlProductsInfo()

	// repo := dao.ProductTypeRepo{}
	// repo.SelectAll()
	// for _, productType := range repo.ProductTypes {
	// 	fmt.Println(productType.Name)
	// }

	// doc := GetDoc(fmt.Sprintf(_productListUrl, 1))

	// doc.Find(`a`).Each(func(i int, s *goquery.Selection) {
	// 	attrKey := `href`
	// 	val, ok := s.Attr(attrKey)
	// 	fmt.Println("===run here 123", val, ok)
	// })

	//TODO: add type here
	// GetProductsInfo(1)
	// GetDetailProductInfoByEndpoint(`/exchanges/poloniex`)
	// GetReviewByEndpoint(`/exchanges/poloniex/review-7bd98a15-8c62-448c-8cab-0f19913f6593`)
}

func GetProductsInfoByPageIndex(pageIndex int) []ProductInfo {
	productInfoList := make([]ProductInfo, 0)
	doc := GetDoc(fmt.Sprintf(_productListUrl, pageIndex))

	DOMKey := `div` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 dzjLTP`)
	doc.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

		DOMKey = `div` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 ReviewTargetCard__Card-sc-qbvmhm-0 jMDOvK kHppZh`)
		s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

			DOMKey = `img`
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `data-src`
				attrVal, foundAttrVal := s.Attr(attrKey)
				fmt.Println("===src: ", attrVal, foundAttrVal)

			})

			DOMKey = `a` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 gtFTPK gOcOhU`)
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

				attrKey := `href`
				attrVal, foundAttrVal := s.Attr(attrKey)
				fmt.Println("===href: ", attrVal, foundAttrVal) //forward to crawl detail
				fmt.Println("===title", s.Text())

			})

			DOMKey = `div` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 bkUBKu`)
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {
				fmt.Println("===type:", s.Text())
			})

			DOMKey = `p` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 ReviewTargetCard__LineClamp-sc-qbvmhm-1 ReviewTargetCard___StyledLineClamp-sc-qbvmhm-2 jVbmuR dUpQvL dxYXvO`)
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {
				fmt.Println("===short description: ", s.Text())
			})

		})

	})

	return productInfoList
}

func GetDetailProductInfoByEndpoint(endpoint string) {
	url := _baseUrl + endpoint
	doc := GetDoc(url)

	DOMKey := `main` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 khpoB`)
	doc.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

		//header(review, name, image ...)
		DOMKey := `div` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Flex-sc-1mngh6p-1 bnwXZr`)
		s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

			DOMKey := `img` + ConvertClassesFormatFromBrowserToGoQuery(`LazyImage-sc-synjzy-0 ReviewTargetLogo__Logo-sc-160quaj-0 HAIBu gWKbJj`)
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {
				val, ok := s.Attr(`data-src`)
				fmt.Println(val, ok)

			})

			DOMKey = `div` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 gmjrOf`)
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {
				fmt.Println(s.Text())
			})

		})

		//footer(review)
		DOMKey = `article` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 Box__Grid-sc-1mngh6p-2 Review__ReviewCard-sc-1xpzhiw-0 iQbXxL cDlhpG`)
		s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

			DOMKey = `a` + ConvertClassesFormatFromBrowserToGoQuery(`Text-sc-kh4piv-0 Anchor-sc-1oa4wrg-0 kKJEOJ dDxbNj`)
			s.Find(DOMKey).Each(func(i int, s *goquery.Selection) {
				attrVal, foundAttr := s.Attr(`href`)
				if foundAttr {
					fmt.Println(`===review`, attrVal)
				}
			})

		})

	})

	//Body (get description, Official website, Social media )
	// DOMKey = `script`
	// doc.Find(DOMKey).Each(func(i int, s *goquery.Selection) {
	// 	val, ok := s.Attr(`type`)
	// 	if ok && val == `application/ld+json` {

	// 		var data any
	// 		json.Unmarshal([]byte(s.Text()), &data)
	// 		description, found := data.(map[string]any)["description"]
	// 		if found {
	// 			fmt.Println(`===description`, description)
	// 		}

	// 		url, found := data.(map[string]any)["url"]
	// 		if found {
	// 			fmt.Println(`===Official website`, url)
	// 		}

	// 		sameAs, found := data.(map[string]any)["sameAs"]
	// 		if found {
	// 			fmt.Println(`===Social media`, sameAs)
	// 		}
	// 	}
	// })

}

func GetReviewByEndpoint(endpoint string) {
	// url := _baseUrl + endpoint
	// doc := GetDoc(url)

	// DOMKey := `main` + ConvertClassesFormatFromBrowserToGoQuery(`Box-sc-1mngh6p-0 khpoB`)
	// doc.Find(DOMKey).Each(func(i int, s *goquery.Selection) {

	// }
}

func ConvertClassesFormatFromBrowserToGoQuery(input string) string {
	classes := input
	classes = `.` + classes
	classes = strings.ReplaceAll(classes, ` `, `.`)
	return classes
}

func GetExchangesName() []string {
	// Load the HTML document
	doc := GetDoc(`https://revain.org/exchanges?sortBy=reviews&page=68`)

	tagAWithClassesForExchangesName := `a.Text-sc-kh4piv-0.Anchor-sc-1oa4wrg-0.gtFTPK.gOcOhU`
	tags := doc.Find(tagAWithClassesForExchangesName)
	exchangesName := make([]string, 0)
	tags.Each(func(i int, tag *goquery.Selection) {
		exchangeName := tag.Text()
		exchangesName = append(exchangesName, exchangeName)
		fmt.Println(exchangeName)
	})

	return exchangesName
}

type ProductInfo struct {
	Image            string
	Title            string
	Type             string
	ShortDescription string
}

type DetailProductInfo struct {
}

func GetListCommentId() {
	// url := `https://revain.org`
}

func GetDetail(exchangeName string) {
	// Load the HTML document
	url := `https://revain.org/exchanges/` + exchangeName
	fmt.Println(url)
	doc := GetDoc(url)

	doc.Find(`img`).Each(func(i int, s *goquery.Selection) {

		fmt.Println(`===== run =====`)

		v, found := s.Attr(`src`)
		if found {
			fmt.Println("attr img ", v)
		} else {
			fmt.Println(`=== not found ===`)
		}

		return
		val, ok := s.Attr("src")
		if ok {
			fmt.Println("==========", val)
		} else {
			fmt.Println("not found src attr", "==========")
		}
	})

	// if ok {
	// 	fmt.Println("======", val)
	// }

	return

	tag := `div.Box-sc-1mngh6p-0.Box__Grid-sc-1mngh6p-2.iNOMXI`
	doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		val := s.Find(`p.Text-sc-kh4piv-0.ReviewContentTeaser__ReviewContent-sc-9r5np5-0.vHQCU.bgSSRv`).First().Text()
		fmt.Println("===========", val)

		val, ok := s.Find(`a.Text-sc-kh4piv-0.Anchor-sc-1oa4wrg-0.kKJEOJ.dDxbNj`).First().Attr(`href`)
		if !ok {
			fmt.Println(`Not found the link`)
		}
		fmt.Println("===========", val)

	})
}
