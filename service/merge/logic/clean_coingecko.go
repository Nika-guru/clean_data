package logic

import (
	"crawler/pkg/log"
	"crawler/service/merge/model/dao"
	dto "crawler/service/merge/model/dto/coingeckgo"
	"fmt"
	"time"
)

const (
	timeScheduleUpdate = 1 * time.Hour
)

func AutoCrawlDataCoinGeckgo() {
	go func() {
		for {
			//######## Crawl coin code and its contract ########
			coinDAOepo := crawTokens()

			for _, coinDAO := range coinDAOepo.Coins {
				product := &dao.Product{}
				sources := map[string]bool{
					`coingecko`: true,
				}
			CheckCoinExist:
				isExist, err := product.IsExist(sources, coinDAO.CoinCode)
				if err != nil {
					log.Println(log.LogLevelError, `service/merge/logic/clean_coingecko.go/AutoCrawlData/product.IsExist(sources)`, err.Error())
					time.Sleep(10 * time.Second)
					goto CheckCoinExist
				}
				if isExist {
					continue
				}

				//############ Crawl more data detail for coin #############
				time.Sleep(5 * time.Second)
				err = crawlDetailByCoinDao(coinDAO)
				if err != nil {
					log.Println(log.LogLevelError, `service/merge/logic/clean_coingecko.go/AutoCrawlData/crawlDetailByCoinDao(coinDAO)`, err.Error())
					continue
				}

				coinDAO.ConvertToProduct(product)
				//test
				err = product.InsertDB()
				if err != nil {
					log.Println(log.LogLevelError, `service/merge/logic/clean_coingecko.go/AutoCrawlData/product.InsertDB()`, err.Error())
				}
			}
			//Wait in duration for next time call API
			time.Sleep(timeScheduleUpdate)
		}

	}()
}

func crawTokens() *dao.CoinRepo {
CrawlTokens:
	DTOrepo := &dto.CoingeckoCoinRepo{}
	err := DTOrepo.CrawlTokens()

	if err != nil {
		log.Println(log.LogLevelError, "crawler/crawTokens", "Crawl function return error: "+err.Error())
		time.Sleep(5 * time.Second) //avoid spam system
		goto CrawlTokens
	}
	log.Println(log.LogLevelInfo, "crawler/crawTokens", "Crawl platform : Successfully")
	coinRepo := &dao.CoinRepo{}
	fmt.Println(`len`, len(DTOrepo.CoingeckoCoins))
	DTOrepo.ConvertTo(coinRepo)
	fmt.Println(`len`, len(coinRepo.Coins))

	return coinRepo
}

func crawlDetailByCoinDao(coinDAO *dao.Coin) error {
	coingeckoDetailCoin := &dto.CoingeckoDetailCoin{}
	coingeckoDetailCoin.CoinCode = coinDAO.CoinCode

	_, err := coingeckoDetailCoin.Crawl() //url, err
	if err != nil {
		return err
	}
	coingeckoDetailCoin.ConvertTo(coinDAO)

	return err
}
