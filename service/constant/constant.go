package constant

import "time"

const (
	MIN_SEC_PER_CALL          time.Duration = 2 * time.Second //Minimum time for a call API
	COINGECKO_TOKENS_API_PATH string        = `https://api.coingecko.com/api/v3/coins/list?include_platform=true`
	MISS_REQUEST_WAIT         time.Duration = 5 * time.Second
	MISS_REQUEST_LIMIT        time.Duration = 30 * time.Minute
)

//crawl data from revain
var ENDPOINTS_PRODUCT_INFO_REVAIN = [...]string{
	`/projects/erc-20`,
	`/projects/trc-20`,
	`/projects/stablecoins`,
	`/projects/defi`,
	`/exchanges`,
	`/wallets`,
	`/blockchain-games`,
	`/crypto-cards`,
	`/mining-pools`,
	`/crypto-trainings`,
	`/categories/nft-marketplaces`,
}

var MAP_CATEGORY_PRODUCT_REVAIN = map[string]bool{
	DEFAULT_CATEGORY_PRODUCT_REVAIN: true,
	`Online Reputation Management`:  true,
	`Crypto Exchanges`:              true,
	`Crypto Wallets`:                true,
	`Blockchain Games`:              true,
	`NFT Marketplaces`:              true,
	`Crypto Cards`:                  true,
	`Bitcoin mining pools`:          true,
	`Crypto Trainings & Courses`:    true,
}

const DEFAULT_CATEGORY_PRODUCT_REVAIN = `Crypto Projects`
