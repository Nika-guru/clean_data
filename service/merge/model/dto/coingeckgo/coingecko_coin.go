package dto

import (
	"base/pkg/log"
	"base/pkg/utils"
	"base/service/merge/model/dao"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	MIN_SEC_PER_CALL          = 2 * time.Second
	COINGECKO_TOKENS_API_PATH = `https://api.coingecko.com/api/v3/coins/list?include_platform=true`
	MISS_REQUEST_WAIT         = 5 * time.Second
)

type CoingeckoCoinRepo struct {
	CoingeckoCoins []CoingeckoCoin
}

type CoingeckoCoin struct {
	CoinID     string `json:"id"`
	CoinSymbol string `json:"symbol"`
	CoinName   string `json:"name"`
	Platform   any    `json:"platforms"`
}

type Contract struct {
	ChainId   string `json:"chainId"`
	ChainName string `json:"chainName"`
	Address   string `json:"address"`
}

// Crawl coin information from external API at [api.coingecko.com].
func (repo *CoingeckoCoinRepo) CrawlTokens() error {
	for {
		time.Sleep(MIN_SEC_PER_CALL)

		coinService := http.Client{
			// Setup 10s Timeout network
			Timeout: 10 * time.Second,
		}

		resp, err := coinService.Get(COINGECKO_TOKENS_API_PATH)

		//Call API error, or timeout
		if err != nil {
			// Continue sleep : Wait to next call API
			time.Sleep(MISS_REQUEST_WAIT)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(log.LogLevelError, "Coingecko/CrawlPlatform", "Convert data to byte : Error")
			return err
		}

		rawCoingeckoDTOs := make([]any, 0)
		err = json.Unmarshal(body, &rawCoingeckoDTOs)
		if err != nil {
			time.Sleep(MISS_REQUEST_WAIT)
			continue
		}

		// Traverse each json object from response array data got above.
		for _, rawCoingeckoDTO := range rawCoingeckoDTOs {
			coingeckoCoinDTO := &CoingeckoCoin{}
			err = utils.Mapping(rawCoingeckoDTO, coingeckoCoinDTO)
			//Mapping from map[string]any to struct{} DTO failed.
			if err != nil {
				return err
			}
			repo.CoingeckoCoins = append(repo.CoingeckoCoins, *coingeckoCoinDTO)
		}
		return nil
	}
}

func (repo *CoingeckoCoinRepo) ConvertTo(DAORepo *dao.CoinRepo) {
	//########## Traverse each DTO ##########
	for _, coingeckoCoin := range repo.CoingeckoCoins {
		coingeckoCoinVal := coingeckoCoin
		platforms := (coingeckoCoin.Platform).(map[string]any)

		coinDAO := &dao.Coin{}
		coinDAO.Src = `coingecko`
		DAORepo.Coins = append(DAORepo.Coins, coinDAO)
		coinDAO.CoinCode = coingeckoCoinVal.CoinID
		coinDAO.CoinSymbol = coingeckoCoinVal.CoinSymbol
		coinDAO.CoinName = coingeckoCoinVal.CoinName

		//########## Native coin ##########
		if len(platforms) == 0 {
			coinDAO.Type = `coin`
			//ChainName, TokenAddress will be NULL
			coinDAO.ChainName = `NULL`
			coinDAO.TokenAddress = `NULL`
		} else
		//########## Token ##########
		{
			coinDAO.Type = `token`

			coinDAO.Contract = make(map[string]any, 0)
			coinDAO.Contract[`contract`] = make([]Contract, 0)
			//########## Traverse each crawled platform from DTO ##########
			idx := 0
			for chainName, tokenAddess := range platforms {
				//###########################  ChainName #########################
				chainNameVal := ``
				//Convert: Special common case
				if chainName == `binance-smart-chain` {
					chainNameVal = `binance`
				} else {
					chainNameVal = chainName
				}

				//###########################  Address #########################
				tokenAddessVal := ``
				if tokenAddess == nil {
					tokenAddessVal = `NULL`
				} else {
					tokenAddessVal = tokenAddess.(string)
				}

				chainList := dao.ChainList{}
				chainList.ChainName = chainNameVal
				chainList.SelectChainIdByChainName()
				//###########################  ChainId #########################
				chainIdVal := ``
				if chainList.ChainId == `` {
					chainIdVal = `NULL` //Non-EVM, or unkown chain
				} else {
					chainIdVal = chainList.ChainId
				}

				//default for value
				if idx == 0 {
					//########## Create new variable, avoid many pointer point last [chainName] after end this loop ##########
					coinDAO.ChainName = chainNameVal
					coinDAO.TokenAddress = tokenAddessVal
					coinDAO.ChainId = chainIdVal
				}

				coinDAO.Contract[`contract`] = append(coinDAO.Contract[`contract`].([]Contract), Contract{
					ChainId:   chainIdVal,
					ChainName: chainNameVal,
					Address:   tokenAddessVal,
				})

				idx++
			}
		}

	}
}

func (repo *CoingeckoCoinRepo) Reset() {
	repo.CoingeckoCoins = make([]CoingeckoCoin, 0)
}
