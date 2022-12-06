package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"review-service/pkg/utils"
	crawl_dappradar_dapp "review-service/service/review/crawl/dappradar_dapp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	_baseUrl        = `https://revain.org`
	_productListUrl = (_baseUrl + `/exchanges?sortBy=reviews&page=%d`)
)

func main() {
	fmt.Println(`start`)
	dom := utils.GetHtmlDomJsRenderByUrl(`https://dappradar.com/ethereum/marketplaces/opensea`)
	crawl_dappradar_dapp.CrawlDappDetailByHtmlDom(dom)
	fmt.Println(`end`)
}

type ReviewReplyRepo struct {
	Total                    uint64        `json:"total"`
	Comments                 []ReviewReply `json:"comments"`
	HasMore                  bool          `json:"hasMore"`
	TotalCommentsWithReplies uint64        `json:"totalCommentsWithReplies"`
}

type ReviewReply struct {
	Author    Author     `json:"author"`
	Content   string     `json:"text"`
	Reactions []Reaction `json:"reactions"`
}

type Author struct {
	FullName string `json:"fullName"`
}

type Reaction struct {
	Account  string `json:"account"`
	Reaction string `json:"reaction"`
}

func testCrawlReply() error {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(`https://revain.org/api/comments?review=8fc7d4bf-1d1d-4f06-95c5-9e48c732d10a&offset=10`) //auto limit 10
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reviewReplyRepo := ReviewReplyRepo{}
	err = json.Unmarshal(body, &reviewReplyRepo)
	fmt.Println(`author: `, reviewReplyRepo.Comments[0].Author.FullName)
	fmt.Println(`content: `, reviewReplyRepo.Comments[0].Content)

	fmt.Println(`content: `, reviewReplyRepo.Comments[0].Reactions[0].Account)
	fmt.Println(`content: `, reviewReplyRepo.Comments[0].Reactions[0].Reaction)

	return err
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
