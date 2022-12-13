package crawl_dappradar_dapp

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//crawl javascript render
//ref: https://www.devdungeon.com/content/web-scraping-go
//lib: https://github.com/geziyor/geziyor

func Crawl() {
	url := fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, constant.ENDPOINT_CHAIN_ALL_DAPPRADAR)
CallAPIBlockchain:
	dom := utils.GetHtmlDomJsRenderByUrl(url)

	if dom == nil {
		time.Sleep(5 * time.Second)
		goto CallAPIBlockchain
	}
	endpointBlockChains := CrawlEndpointBlockchainsByHtmlDom(dom)
	if len(endpointBlockChains) == 0 {
		time.Sleep(5 * time.Second)
		goto CallAPIBlockchain
	}

	dtoEndpointBlockchainRepo := dto_dappradar.EndpointBlockchainRepo{}
	dtoEndpointBlockchainRepo.EndpointBlockchains = endpointBlockChains
	daoBlockchainRepo := &dao.BlockchainRepo{}
	dtoEndpointBlockchainRepo.ConvertTo(daoBlockchainRepo)
	daoBlockchainRepo.InsertDB()

	for _, endpointBlockChain := range endpointBlockChains {
		for pageIdx := 1; ; pageIdx++ {
			url := fmt.Sprintf(`%s%v/%d`, constant.BASE_URL_DAPPRADAR, endpointBlockChain.Endpoint, pageIdx)
		CallAPIListPagination:
			dom := utils.GetHtmlDomJsRenderByUrl(url)
			if dom == nil {
				log.Println(log.LogLevelDebug, `service/review/crawl/dappradar_dapp/Crawl/getDomJsLoad`, `dom get by js loading is nil`)
				pageIdx--
				time.Sleep(5 * time.Second)
				continue
			}
			if IsEndPage(dom) {
				break
			}

			endpointDappList := CrawlDappEndpointByHtmlDom(dom)
			//Response without data
			if len(endpointDappList) == 0 {
				time.Sleep(5 * time.Second)
				goto CallAPIListPagination
			}

			for _, endpointDapp := range endpointDappList {
			CallAPIDetailPage:
				url := fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, endpointDapp.Endpoint)
				dom := utils.GetHtmlDomJsRenderByUrl(url)

				//net::ERR_INTERNET_DISCONNECTED
				if dom == nil {
					time.Sleep(10 * time.Second)
					continue
				}
				isPasss := CrawlDappDetailByHtmlDom(dom, endpointDapp.DetailDapp)
				if !isPasss {
					time.Sleep(5 * time.Second)
					fmt.Println(log.LogLevelError, `dom fail `+url)
					goto CallAPIDetailPage
				}
				err := endpointDapp.InsertDB()
				if err != nil {
					fmt.Println(log.LogLevelError, `Insert fail `+url, err.Error())
					continue
				}
			}
		}
	}
}

func IsEndPage(dom *goquery.Document) bool {
	isEndPage := false
	domKey := `h2`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		_txtNotifyEndPageDappRadar := `Please change the filters to explore more`
		if s.Text() == _txtNotifyEndPageDappRadar {
			isEndPage = true
		}
	})
	return isEndPage
}
