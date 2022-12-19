package dto

import (
	"base/pkg/log"
	"base/service/merge/model/dao"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	detailCoinUrl string = `https://api.coingecko.com/api/v3/coins/%s`
)

type CoingeckoDetailCoin struct {
	CoinCode        string         //must inititalize in constructor
	Descriptions    any            `json:"description"`
	Images          any            `json:"image"`
	Links           LinkDetailCoin `json:"links"`
	Platforms       any            `json:"platforms"`
	DetailPlatforms any            `json:"detail_platforms"` //k: chainName; v:{decimal_place: any.(int); contract_address:any.(string)}
	Categories      []string       `json:"categories"`
	MarketData      MarketData     `json:"market_data"`

	//Add extra to convert dto
	ChainName    *string `json:"chainName"`
	TokenAddress any     `json:"tokenAddress"`
}
type MarketData struct {
	TotalSupply       any `json:"total_supply"`
	MaxSupply         any `json:"max_supply"`
	CirculatingSupply any `json:"circulating_supply"`
}

type LinkDetailCoin struct {
	Homepage                    []string                `json:"homepage"`
	BlockchainSite              []string                `json:"blockchain_site"`
	OfficialForumUrl            []string                `json:"official_forum_url"`
	ChatUrl                     []string                `json:"chat_url"`
	AnnouncementUrl             []string                `json:"announcement_url"`
	TwitterScreenName           string                  `json:"twitter_screen_name"`
	FacebookUsername            string                  `json:"facebook_username"`
	BitcointalkThreadIdentifier *int                    `json:"bitcointalk_thread_identifier"`
	TelegramChannelIdentifier   string                  `json:"telegram_channel_identifier"`
	SubredditUrl                string                  `json:"subreddit_url"`
	ReposUrl                    CoingeckoDetailReposUrl `json:"repos_url"`
}

type CoingeckoDetailReposUrl struct {
	Github    []string `json:"github"`
	Bitbucket []string `json:"bitbucket"`
}

type CoingeckoDetailContract struct {
	DecimalPlace    *int   `json:"decimal_place"`
	ContractAddress string `json:"contract_address"` //absent data ""
}

type FailDTO struct {
	Status struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_message"`
	} `json:"status"`
}

type NotExistedModel struct {
	Error string `json:"error"`
}

func (coingeckoDTO *CoingeckoDetailCoin) Crawl() (url string, err error) {
	for {
		time.Sleep(5 * time.Second)

		client := http.Client{
			Timeout: 10 * time.Second,
		}
		url := fmt.Sprintf(detailCoinUrl, coingeckoDTO.CoinCode)
		resp, err := client.Get(url)
		if err != nil { // Call API from coingecko error, or timeout, network missing.
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "API from coingecko error, or timeout, network missing. DETAIL: "+err.Error())
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil { // Data got, is not valid to convert to data type []byte.
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Conver data type error. DETAIL: "+err.Error())
			continue
		}

		notExistedModel := &NotExistedModel{}
		err = json.Unmarshal(body, notExistedModel)
		if err != nil {
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Unmarshal data error. DETAIL: "+err.Error())
			continue
		}
		//Don't existed coin
		if notExistedModel.Error == "coin not found" {
			// fmt.Println("notExistedModel", notExistedModel, notExistedModel.Error)
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Not found coin from API(coin deleted) ")
			return url, errors.New("not found coin from API(coin deleted)")
		}

		failDTO := &FailDTO{}
		err = json.Unmarshal(body, failDTO)
		if err != nil {
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Unmarshal FailDTO error "+err.Error())
			continue
		}

		//Rate limit
		if failDTO.Status.ErrorCode != 429 {

			err = json.Unmarshal(body, coingeckoDTO)
			if err != nil { // Unzip data got as datatype []byte to raw DTO struct failed.
				log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Unmarshal coingeckoDTO error  "+err.Error())
				continue
			} else { //ko loi
				return url, nil
			}

		} else { //Bi chan
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Block due to time limit ")
			time.Sleep(10 * time.Second)
			continue
		}

	}
}

func (dto *CoingeckoDetailCoin) ConvertTo(dao *dao.Coin) {
	//################################ Detail JSONB #################################
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	//Fixed
	dao.Detail[`productCodeCoingecko`] = dao.CoinCode
	dao.Detail[`maxSupply`] = dto.MarketData.MaxSupply
	dao.Detail[`decimals`] = dto.DetailPlatforms

	//Updated
	dao.Detail[`totalSupply`] = dto.MarketData.CirculatingSupply
	dao.Detail[`marketcap`] = 0.0
	dao.Detail[`volumeTrading`] = 0.0
	dao.Detail[`holder`] = 0.0

	//1######### Update description of native/token coin #########
	for locale, description := range (dto.Descriptions).(map[string]any) {
		//get only description english
		if locale == "en" {
			if description != nil {
				descriptionVal := description.(string)
				dao.Description = descriptionVal
			} else {
				dao.Description = ``
			}
			break
		}
	}

	//2######### Update Image of native/token coin #########
	//3######### prepare update image thumb of LinkIcon( contract link item with must be native coin) #########
	for imgSizeType, imgUrl := range (dto.Images).(map[string]any) {
		//Update coins tbl: (native/ token) info
		if imgSizeType == "large" {
			if imgUrl != nil {
				imgVal := imgUrl.(string)
				dao.CoinImage = imgVal
			}
			break
		}
	}

	var tokenDecimal any = nil
	//4######### Update decimal of token coin #########
	for chainName, coingeckoDetailContract := range (dto.DetailPlatforms).(map[string]any) {
		//Token coin --> same chain name + same coin id in url
		if chainName != "" && dto.ChainName != nil && chainName == *dto.ChainName {
			coingeckoDetailContractVal := coingeckoDetailContract
			//check match cointract(chainname + token address)

			for key, value := range coingeckoDetailContractVal.(map[string]any) {
				if key == "decimal_place" {
					tokenDecimal = value
				}
			}

		}
	}
	if tokenDecimal != nil {
		val := tokenDecimal.(float64)
		dao.Decimals = uint8(val)
	} else { //Native coin, or eosblack token, where contract address and decimal is nil
		dao.Decimals = 0
	}

	// ############### homepage ###############
	urls := make([]string, 0)
	for _, url := range dto.Links.Homepage {
		urlVal := url
		//don't append blank link
		if urlVal != "" {
			urls = append(urls, urlVal)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found := dao.Detail[`website`]
	if !found {
		website := make(map[string]any, 0)
		website[`homepage`] = urls
		dao.Detail[`website`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`homepage`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`homepage`] = urls
		} else {
			website = append(oldWebsite.([]string), urls...)
		}
		dao.Detail[`website`] = website
	}

	// ############### Announcement ###############
	announcementUrls := make([]string, 0)
	for _, url := range dto.Links.AnnouncementUrl {
		urlVal := url
		//don't append blank link
		if url != "" {
			announcementUrls = append(announcementUrls, urlVal)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found = dao.Detail[`website`]
	if !found {
		website := make(map[string]any, 0)
		website[`announcement`] = announcementUrls
		dao.Detail[`website`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`announcement`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`announcement`] = announcementUrls
		} else {
			website = append(oldWebsite.([]string), announcementUrls...)
		}
		dao.Detail[`website`] = website
	}

	// ############### Blockchain site ###############
	blockchainSiteUrls := make([]string, 0)
	for _, url := range dto.Links.BlockchainSite {
		urlVal := url
		//don't append blank link
		if url != "" {
			blockchainSiteUrls = append(blockchainSiteUrls, urlVal)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found = dao.Detail[`website`]
	if !found {
		website := make(map[string]any, 0)
		website[`blockchainSite`] = blockchainSiteUrls
		dao.Detail[`website`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`blockchainSite`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`blockchainSite`] = blockchainSiteUrls
		} else {
			website = append(oldWebsite.([]string), blockchainSiteUrls...)
		}
		dao.Detail[`website`] = website
	}

	// ############### Subreddit ###############
	if dto.Links.SubredditUrl != "" {
		subredditUrl := dto.Links.SubredditUrl
		subredditUrls := make([]string, 0)
		subredditUrls = append(subredditUrls, subredditUrl)
		if dao.Detail == nil {
			dao.Detail = make(map[string]any)
		}
		website, found = dao.Detail[`community`]
		if !found {
			website := make(map[string]any, 0)
			website[`subreddit`] = subredditUrls
			dao.Detail[`community`] = website
		} else {
			oldWebsite, found := website.(map[string]any)[`subreddit`]
			//append new to existed map
			if !found {
				website.(map[string]any)[`subreddit`] = subredditUrls
			} else {
				website = append(oldWebsite.([]string), subredditUrls...)
			}
			dao.Detail[`community`] = website
		}
	}

	// ############### Twitter ###############
	if dto.Links.TwitterScreenName != "" {
		twitterUrl := fmt.Sprintf("https://twitter.com/%s", dto.Links.TwitterScreenName)

		twitterUrls := make([]string, 0)
		twitterUrls = append(twitterUrls, twitterUrl)
		if dao.Detail == nil {
			dao.Detail = make(map[string]any)
		}
		website, found = dao.Detail[`community`]
		if !found {
			website := make(map[string]any, 0)
			website[`twitter`] = twitterUrls
			dao.Detail[`community`] = website
		} else {
			oldWebsite, found := website.(map[string]any)[`twitter`]
			//append new to existed map
			if !found {
				website.(map[string]any)[`twitter`] = twitterUrls
			} else {
				website = append(oldWebsite.([]string), twitterUrls...)
			}
			dao.Detail[`community`] = website
		}
	}

	// ############### Facebook ###############
	if dto.Links.FacebookUsername != "" {
		facebookUrl := fmt.Sprintf("https://www.facebook.com/%s", dto.Links.FacebookUsername)
		facebookUrls := make([]string, 0)
		facebookUrls = append(facebookUrls, facebookUrl)
		if dao.Detail == nil {
			dao.Detail = make(map[string]any)
		}
		website, found = dao.Detail[`community`]
		if !found {
			website := make(map[string]any, 0)
			website[`facebook`] = facebookUrls
			dao.Detail[`community`] = website
		} else {
			oldWebsite, found := website.(map[string]any)[`facebook`]
			//append new to existed map
			if !found {
				website.(map[string]any)[`facebook`] = facebookUrls
			} else {
				website = append(oldWebsite.([]string), facebookUrls...)
			}
			dao.Detail[`community`] = website
		}
	}

	// ############### Chat url ###############
	chatUrls := make([]string, 0)
	for _, url := range dto.Links.ChatUrl {
		//don't append blank link
		if url != "" {
			chatUrls = append(chatUrls, url)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found = dao.Detail[`community`]
	if !found {
		website := make(map[string]any, 0)
		website[`chatUrl`] = chatUrls
		dao.Detail[`community`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`chatUrl`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`chatUrl`] = chatUrls
		} else {
			website = append(oldWebsite.([]string), chatUrls...)
		}
		dao.Detail[`community`] = website
	}

	// ############### Official Forum ###############
	officialForumUrls := make([]string, 0)
	for _, url := range dto.Links.OfficialForumUrl {
		//don't append blank link
		if url != "" {
			urlVal := url
			officialForumUrls = append(officialForumUrls, urlVal)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found = dao.Detail[`community`]
	if !found {
		website := make(map[string]any, 0)
		website[`officialForum`] = officialForumUrls
		dao.Detail[`community`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`officialForum`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`officialForum`] = officialForumUrls
		} else {
			website = append(oldWebsite.([]string), officialForumUrls...)
		}
		dao.Detail[`community`] = website
	}

	// ############### Bitcointalk ###############
	if dto.Links.BitcointalkThreadIdentifier != nil {
		bitcointalkUrl := fmt.Sprintf("https://bitcointalk.org/index.php?topic=%d", *dto.Links.BitcointalkThreadIdentifier)

		bitcointalkUrls := make([]string, 0)
		bitcointalkUrls = append(bitcointalkUrls, bitcointalkUrl)
		if dao.Detail == nil {
			dao.Detail = make(map[string]any)
		}
		website, found = dao.Detail[`community`]
		if !found {
			website := make(map[string]any, 0)
			website[`bitcointalk`] = bitcointalkUrls
			dao.Detail[`community`] = website
		} else {
			oldWebsite, found := website.(map[string]any)[`bitcointalk`]
			//append new to existed map
			if !found {
				website.(map[string]any)[`bitcointalk`] = bitcointalkUrls
			} else {
				website = append(oldWebsite.([]string), bitcointalkUrls...)
			}
			dao.Detail[`community`] = website
		}
	}

	// ############### Telegram ###############
	if dto.Links.TelegramChannelIdentifier != "" {
		telegramUrl := fmt.Sprintf("https://t.me/%s", dto.Links.TelegramChannelIdentifier)
		telegramUrls := make([]string, 0)
		telegramUrls = append(telegramUrls, telegramUrl)
		if dao.Detail == nil {
			dao.Detail = make(map[string]any)
		}
		website, found = dao.Detail[`community`]
		if !found {
			website := make(map[string]any, 0)
			website[`telegram`] = telegramUrls
			dao.Detail[`community`] = website
		} else {
			oldWebsite, found := website.(map[string]any)[`telegram`]
			//append new to existed map
			if !found {
				website.(map[string]any)[`telegram`] = telegramUrls
			} else {
				website = append(oldWebsite.([]string), telegramUrls...)
			}
			dao.Detail[`community`] = website
		}
	}

	// ############### Twitter search ###############
	if dao.CoinSymbol != `` {
		twitterSearchUrl := fmt.Sprintf(`https://twitter.com/search?q=$%s`, dao.CoinSymbol)
		twitterSearchUrls := make([]string, 0)
		twitterSearchUrls = append(twitterSearchUrls, twitterSearchUrl)
		if dao.Detail == nil {
			dao.Detail = make(map[string]any)
		}
		website, found = dao.Detail[`community`]
		if !found {
			website := make(map[string]any, 0)
			website[`twitterSearch`] = twitterSearchUrls
			dao.Detail[`community`] = website
		} else {
			oldWebsite, found := website.(map[string]any)[`twitterSearch`]
			//append new to existed map
			if !found {
				website.(map[string]any)[`twitterSearch`] = twitterSearchUrls
			} else {
				website = append(oldWebsite.([]string), twitterSearchUrls...)
			}
			dao.Detail[`community`] = website
		}
	}

	// ############### Github ###############
	githubUrls := make([]string, 0)
	for _, url := range dto.Links.ReposUrl.Github {
		//don't append blank link
		if url != "" {
			githubUrls = append(githubUrls, url)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found = dao.Detail[`sourceCode`]
	if !found {
		website := make(map[string]any, 0)
		website[`github`] = githubUrls
		dao.Detail[`sourceCode`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`github`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`github`] = githubUrls
		} else {
			website = append(oldWebsite.([]string), githubUrls...)
		}
		dao.Detail[`sourceCode`] = website
	}

	// ############### Bitbucket ###############
	bitbucketUrls := make([]string, 0)
	for _, url := range dto.Links.ReposUrl.Bitbucket {
		//don't append blank link
		if url != "" {
			bitbucketUrls = append(bitbucketUrls, url)
		}
	}
	if dao.Detail == nil {
		dao.Detail = make(map[string]any)
	}
	website, found = dao.Detail[`sourceCode`]
	if !found {
		website := make(map[string]any, 0)
		website[`bitbucket`] = bitbucketUrls
		dao.Detail[`sourceCode`] = website
	} else {
		oldWebsite, found := website.(map[string]any)[`bitbucket`]
		//append new to existed map
		if !found {
			website.(map[string]any)[`bitbucket`] = bitbucketUrls
		} else {
			website = append(oldWebsite.([]string), bitbucketUrls...)
		}
		dao.Detail[`sourceCode`] = website
	}

	// ############### Tag --> Category ###############
	categories := ``
	for idx, category := range dto.Categories {
		categories += category
		if idx != len(dto.Categories)-1 {
			categories += `,`
		}
	}
	dao.Tag = categories
}
