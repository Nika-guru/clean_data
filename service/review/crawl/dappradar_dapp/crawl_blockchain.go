package crawl_dappradar_dapp

import (
	"review-service/pkg/utils"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlEndpointBlockchainsByHtmlDom(dom *goquery.Document) []*dto_dappradar.EndpointBlockchain {
	blockchains := make([]*dto_dappradar.EndpointBlockchain, 0)

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`rankings-filters`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`rankings-filters__protocols`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			existedBlockchainEndpoint := make(map[string]bool)
			domKey := `a` + utils.ConvertClassesFormatFromBrowserToGoQuery(`rankings-filters__protocol`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				blockchain := &dto_dappradar.EndpointBlockchain{}

				endpointBlockchain, foundEndpointBlockchain := s.Attr(`href`)
				if foundEndpointBlockchain {
					blockchainName := s.Text()
					_, foundBlockchainEndpoint := existedBlockchainEndpoint[endpointBlockchain]
					if !foundBlockchainEndpoint {
						existedBlockchainEndpoint[endpointBlockchain] = true

						_allBlockchain := `All`
						_endpointAllBlockchain := `/rankings`
						//don't get endpoint all
						if endpointBlockchain != _endpointAllBlockchain && blockchainName != _allBlockchain {
							blockchain.BlockchainName = blockchainName
							blockchain.Endpoint = endpointBlockchain

							blockchainId := strings.ReplaceAll(endpointBlockchain, `/rankings/protocol/`, ``)
							blockchain.BlockchainId = blockchainId

							_otherBlockchain := `Other`
							//Identified blockchain
							if blockchainName != _otherBlockchain {
								htmlTxt, err := s.Html()
								if err == nil {
									neededHtml := strings.Split(htmlTxt, `</svg>`)[0]
									blockchain.BlockchainImageSvg = neededHtml
								}
							} else
							//Other blockchain
							{
								_otherBlockchainImgSvg := `<svg width="12px" height="12px" viewBox="0 0 17 16" xmlns="http://www.w3.org/2000/svg" fill="#B1BBCE" mr="4px" class="sc-eBOGjE fXoYMz"><path d="M15.174 7.475c.759 0 1.375-.617 1.375-1.375V1.875c0-.758-.616-1.375-1.375-1.375H10.95c-.758 0-1.375.617-1.375 1.375V3.15H7.641V1.875C7.641 1.117 7.024.5 6.266.5H2.041C1.283.5.666 1.117.666 1.875V6.1c0 .758.617 1.375 1.375 1.375h1.275v1.383H2.041c-.758 0-1.375.617-1.375 1.375v4.225c0 .759.617 1.375 1.375 1.375h4.225c.758 0 1.375-.616 1.375-1.375v-1.283h1.933v1.283c0 .759.617 1.375 1.375 1.375h4.225c.759 0 1.375-.616 1.375-1.375v-4.225c0-.758-.616-1.375-1.375-1.375H13.9V7.475h1.275Zm-3.933-5.308h3.642v3.641H11.24V2.167Zm-8.908 0h3.641v3.641H2.333V2.167Zm3.641 12H2.333v-3.642h3.641v3.642Zm8.909 0H11.24v-3.642h3.642v3.642Zm-2.659-5.309h-1.283c-.758 0-1.375.617-1.375 1.375v1.275H7.633v-1.275c0-.758-.617-1.375-1.375-1.375H4.974V7.475h1.284c.758 0 1.375-.617 1.375-1.375V4.817h1.933V6.1c0 .758.617 1.375 1.375 1.375h1.283v1.383Z"></path></svg>`
								blockchain.BlockchainImageSvg = _otherBlockchainImgSvg
							}
							blockchains = append(blockchains, blockchain)
						}
					}
				}
			})

		})

	})

	return blockchains
}
