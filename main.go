package main

import (
	"encoding/json"
	"fmt"
	"log"
	"review-service/pkg/utils"
	crawl "review-service/service/review/crawl/revain"

	"github.com/PuerkitoBio/goquery"
)

const (
	_baseUrl        = `https://revain.org`
	_productListUrl = (_baseUrl + `/exchanges?sortBy=reviews&page=%d`)
)

func main() {
	dom, err := crawl.GetHtmlDomByUrl(`https://revain.org/projects/revain/review-8fc7d4bf-1d1d-4f06-95c5-9e48c732d10a`)
	if err != nil {
		log.Fatal(err)
	}

	extract(dom)
}

func extract(dom *goquery.Document) {
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

			// positiveNotes, found := data.(map[string]any)[`positiveNotes`]
			// if found {

			// 	for _, val := range positiveNotes.(map[string]any)[`itemListElement`].([]any) {
			// 		fmt.Println(`======itemListElement`, val.(map[string]any)[`name`])
			// 	}
			// }

			// negativeNotes, found := data.(map[string]any)[`negativeNotes`]
			// if found {
			// 	for _, val := range negativeNotes.(map[string]any)[`itemListElement`].([]any) {
			// 		fmt.Println(`======itemListElement`, val.(map[string]any)[`name`])
			// 	}
			// }

			// author, found := data.(map[string]any)[`author`]
			// if found {
			// 	fmt.Println(`author name`, author.([]any)[0].(map[string]any)[`name`])
			// 	fmt.Println(`author name`, author.([]any)[0].(map[string]any)[`image`])
			// }

			reviewRating, found := data.(map[string]any)[`reviewRating`]
			if found {
				fmt.Println(`ratingValue`, reviewRating.(map[string]any)[`ratingValue`])
			}

			datePublished, found := data.(map[string]any)[`datePublished`]
			if found {
				time, err := utils.StringToTimestamp(datePublished.(string))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(`datePublished`, datePublished, time)
			}

		}
	})
}
