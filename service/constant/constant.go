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

const (
	RESP_SUCCESS_STATUS_CODE      = 200
	RESP_TOO_MANY_REQ_STATUS_CODE = 429
	RESP_NOT_FOUND_STATUS_CODE    = 404
	WAIT_DURATION_WHEN_RATE_LIMIT = 5 * time.Second
)

//Cache
const (
	KEY_CACHE_REVAIN_PRODUCT_INFO             = `KEY_CACHE_REVAIN_PRODUCT_INFO`
	KEY_CACHE_COINGECKO_PRODUCT_CATEGORY_INFO = `KEY_CACHE_COINGECKO_PRODUCT_CATEGORY_INFO`
)

// dappradar
const (
	BASE_URL_DAPPRADAR           = `https://dappradar.com`
	ENDPOINT_CHAIN_ALL_DAPPRADAR = `/rankings`
)

// top100token
const (
	BASE_URL_TOP100TOKEN        = `https://top100token.com`
	ENDPOINT_LATEST_TOP100TOKEN = `/latest`
)

// coincarp
const (
	URL_GET_INVESTORS_COINCARP      = `https://sapi.coincarp.com/api/v1/market/investor/list?start=%d&length=%d`
	PARAM_SORTING_INVESTOR_COINCARP = "&lang=en-US&draw=1&columns%5B0%5D%5Bdata%5D=investorname&columns%5B0%5D%5Bname%5D=&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=categoryname&columns%5B1%5D%5Bname%5D=&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=false&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=location&columns%5B2%5D%5Bname%5D=&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=false&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=launched&columns%5B3%5D%5Bname%5D=&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=false&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=projectlist&columns%5B4%5D%5Bname%5D=&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=true&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=investorcount&columns%5B5%5D%5Bname%5D=&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=false&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&order%5B0%5D%5Bcolumn%5D=4&order%5B0%5D%5Bdir%5D=desc&search%5Bvalue%5D=&search%5Bregex%5D=false&_=1670742863427"
	LIMIT_REC_INVESTOR_COINCARP     = 20

	//SORT BY FUND DATE LATEST
	URL_GET_FUNDRAISINGS_COINCARP = `https://sapi.coincarp.com/api/v1/market/fundraising/list?start=%d&length=%d`
	// PARAMS_SORTING_FUNDRASINGS_COINCARP = "&lang=en-US&draw=1&columns%5B0%5D%5Bdata%5D=projectname&columns%5B0%5D%5Bname%5D=&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=categorylist&columns%5B1%5D%5Bname%5D=&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=false&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=fundstagename&columns%5B2%5D%5Bname%5D=&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=false&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=fundamount&columns%5B3%5D%5Bname%5D=&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=investorlist&columns%5B4%5D%5Bname%5D=&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=false&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=funddate&columns%5B5%5D%5Bname%5D=&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=true&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&order%5B0%5D%5Bcolumn%5D=5&order%5B0%5D%5Bdir%5D=desc&search%5Bvalue%5D=&search%5Bregex%5D=false&_=1670814896415"
	LIMIT_REC_FUND_RAISINGS_COINCARP = 20

	URL_GET_INVESTOR_BY_FUND_RAISING_COINCARP = `https://sapi.coincarp.com/api/v1/market/project/investorlist?projectcode=%s&fundcode=%s&page=%d&pagesize=%d&lang=en-US`
	LIMIT_INVESTOR_BY_FUND_RAISING_COINCARP   = 20

	BASE_URL_COINCARP              = `https://www.coincarp.com`
	ENDPOINT_FUND_RAISING_COINCARP = `/fundraising`
	ENDPOINT_INVESTOR_COINCARP     = `/investor`
	BASE_URL_IMAGE_SERVER_COINCARP = `https://s1.coincarp.com`
)
