package dto

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	"time"
)

type CoinRepo struct {
	Coins []Coin
}

type Coin struct {
	CoinId     string `json:"id"`
	CoinSymbol string `json:"symbol"`
	CoinName   string `json:"name"`
	Platform   any    `json:"platforms"`
}

// Crawl coin information from external API at [api.coingecko.com].
func (repo *CoinRepo) CrawlCoins() error {
	missRequest := 0

	for {
		time.Sleep(constant.MIN_SEC_PER_CALL)

		coinService := http.Client{
			// Setup 10s Timeout network
			Timeout: 10 * time.Second,
		}

		resp, err := coinService.Get(constant.COINGECKO_TOKENS_API_PATH)

		//Call API error, or timeout
		if err != nil {
			// Continue sleep : Wait to next call API
			time.Sleep(constant.MISS_REQUEST_WAIT)

			// increment miss request
			missRequest += 1
			if missRequest >= int(constant.MISS_REQUEST_LIMIT/constant.MISS_REQUEST_WAIT) {
				log.Println(log.LogLevelWarn, "Coingecko/CrawlPlatform", "Get 3rd API Coingecko Timeout")
				return err
			}
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(log.LogLevelWarn, "Coingecko/CrawlPlatform", "Convert data to byte : Error")
			return err
		}

		rawCoingeckoDTOs := make([]any, 0)
		err = json.Unmarshal(body, &rawCoingeckoDTOs)
		if err != nil {
			log.Println(log.LogLevelWarn, "Coingecko/CrawlPlatform", "3rd API coingecko block due to over rate limit, detail: Unmarshal error "+err.Error())

			// increment miss request
			missRequest += 1
			if missRequest >= int(constant.MISS_REQUEST_LIMIT/constant.MISS_REQUEST_WAIT) {
				log.Println(log.LogLevelWarn, "Coingecko/CrawlPlatform", "Get 3rd API Coingecko Failed Because Rate Limit exceeded")
				return err
			}
			continue
		}

		// Traverse each json object from response array data got above.
		for _, rawCoingeckoDTO := range rawCoingeckoDTOs {
			coinDTO := &Coin{}
			err = utils.Mapping(rawCoingeckoDTO, coinDTO)
			//Mapping from map[string]any to struct{} DTO failed.
			if err != nil {
				return err
			}
			repo.Coins = append(repo.Coins, *coinDTO)
		}
		return nil
	}
}
