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

const (
	KEY_CACHE_PRODUCT_DETAIL = `KEY_CACHE_PRODUCT_DETAIL`
)
