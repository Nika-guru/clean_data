package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"review-service/pkg/log"
	"review-service/service/coin_info/model/dao"
	"strings"
	"time"
)

const (
	detailCoinUrl    string        = `https://api.coingecko.com/api/v3/coins/%s`
	_minSecPerCall   time.Duration = 2 * time.Second
	_missRequestWait time.Duration = 5 * time.Second
)

//Don't have check exist for link item
func CrawlDetail(insertedCoinIdDetail *dao.CoinRepo) {
	coinRepo := &dao.CoinRepo{}
	linkItemUrlRepo := &dao.LinkItemRepo{}
	existedLinkItemUrlRepo := &dao.LinkItemRepo{}
	existedLinkItemUrlRepo.SelectAll()
	existedLinkItemMap := make(map[string]bool, 0)
	//Traverse map to search
	for _, linkItem := range existedLinkItemUrlRepo.LinkItems {
		_, found := existedLinkItemMap[linkItem.GetKeyMap()]
		if !found {
			existedLinkItemMap[linkItem.GetKeyMap()] = true
		}
	}

	coinTagRepo := &dao.CoinTagRepo{}

	//CoinId + Chainname + TokenAddress
	oldCoinUnique := ""
	oldCoinId := ""
	//for here list coin inserted
	for _, existedCoin := range insertedCoinIdDetail.Coins {
		if existedCoin.CoinID != nil {
			coinRef := existedCoin
			//For update
			time.Sleep(5 * time.Second)
			coingeckoDetailCoin := &CoingeckoDetailCoin{}
			coingeckoDetailCoin.CoinId = *coinRef.CoinID
			err := coingeckoDetailCoin.Crawl()
			if err != nil {
				log.Println(log.LogLevelWarn, "dto/detail_coins/coingecko/CrawlDetail", "this coin id not exist, go next. DETAIL: "+err.Error())
				time.Sleep(10 * time.Second)
				continue
			}

			//New coin unique(CoinId + Chainname + TokenAddress)
			if oldCoinUnique != coinRef.GetKeyMap() {
				coinRepo.Coins = append(coinRepo.Coins, coinRef)

				//1######### Update description of native/token coin #########
				for locale, description := range (coingeckoDetailCoin.Descriptions).(map[string]any) {
					//get only description english
					if locale == "en" {
						if description != nil {
							val := description.(string)
							coinRef.CoinDescription = &val
						}
						break
					}
				}

				//2######### Update Image of native/token coin #########
				//3######### prepare update image thumb of LinkIcon( contract link item with must be native coin) #########
				for imgSizeType, imgUrl := range (coingeckoDetailCoin.Images).(map[string]any) {
					//Update coins tbl: (native/ token) info
					if imgSizeType == "large" {
						if imgUrl != nil {
							val := imgUrl.(string)
							coinRef.CoinImage = &val
						}
						break
					}
				}

				//4######### Update decimal of token coin #########
				for chainName, coingeckoDetailContract := range (coingeckoDetailCoin.DetailPlatforms).(map[string]any) {
					//Token coin
					if chainName != "" {
						//eosblack chainname == eos + contract address == nil --> decimal == nil
						if coinRef.ChainName != nil && coinRef.TokenAddress != nil {
							coingeckoDetailContractVal := coingeckoDetailContract
							//check match cointract(chainname + token address)

							var tokenDecimal any
							for key, value := range coingeckoDetailContractVal.(map[string]any) {
								if key == "decimal_place" {
									tokenDecimal = value
								} else if key == "contract_address" {
									tokenAddress := value

									if *coinRef.ChainName == chainName && coinRef.TokenAddress.(string) == tokenAddress {
										// fmt.Println("1==============chainname", *coinRef.ChainName == chainName, *coinRef.ChainName, chainName)
										// fmt.Println("2==============tokenaddress", coinRef.TokenAddress.(string) == tokenAddress, coinRef.TokenAddress, tokenAddress)
										// fmt.Println("3===============decimal", tokenDecimal)
										if tokenDecimal != nil {
											val := tokenDecimal.(float64)
											coinRef.TokenDecimal = &val
										}
									}
								}

							}

						}
					}
				}

				coinRepo.UpdateDBDecimalImageDescription()
				// updatedCoin := coinRepo.UpdateDBDecimalImageDescription()
				// fmt.Println("1===update success", updatedCoin, "coin/token")
				//for next loop reset
				coinRepo.Coins = make([]*dao.Coin, 0)

				// Update current coin unique to old coin unique for next loop
				oldCoinUnique = coinRef.GetKeyMap()
			}

			//New coin id
			if coinRef.CoinID != nil && oldCoinId != *coinRef.CoinID {
				//######### Start: Insert link item #########
				groupWebsite := 1
				for _, url := range coingeckoDetailCoin.Links.Homepage {
					urlVal := url
					//don't append blank link
					if url != "" {
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupWebsite,
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}
				for _, url := range coingeckoDetailCoin.Links.AnnouncementUrl {
					urlVal := url
					//don't append blank link
					if url != "" {
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupWebsite,
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}

				groupExplorer := 2
				for _, url := range coingeckoDetailCoin.Links.BlockchainSite {
					urlVal := url
					//don't append blank link
					if url != "" {
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupExplorer,
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}

				groupCommunity := 4
				if coingeckoDetailCoin.Links.SubredditUrl != "" {
					urlVal := coingeckoDetailCoin.Links.SubredditUrl
					thumbImgRedit := "fa-brands fa-reddit"
					titleReddit := "Reddit"
					linkItem := &dao.LinkItem{
						LinkUrl:     &urlVal,
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupCommunity,
						LinkIcon:    &thumbImgRedit,
						LinkTitle:   &titleReddit,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				if coingeckoDetailCoin.Links.TwitterScreenName != "" {
					urlVal := fmt.Sprintf("https://twitter.com/%s", coingeckoDetailCoin.Links.TwitterScreenName)
					thumbImgTwitter := "fa-brands fa-twitter"
					titleTwitter := "Twitter"
					linkItem := &dao.LinkItem{
						LinkUrl:     &urlVal,
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupCommunity,
						LinkIcon:    &thumbImgTwitter,
						LinkTitle:   &titleTwitter,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				if coingeckoDetailCoin.Links.FacebookUsername != "" {
					urlVal := fmt.Sprintf("https://www.facebook.com/%s", coingeckoDetailCoin.Links.FacebookUsername)
					thumbImgFacebook := "fa-brands fa-facebook"
					titleFacebook := "Facebook"
					linkItem := &dao.LinkItem{
						LinkUrl:     &urlVal,
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupCommunity,
						LinkIcon:    &thumbImgFacebook,
						LinkTitle:   &titleFacebook,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				for _, url := range coingeckoDetailCoin.Links.ChatUrl {
					//don't append blank link
					if url != "" {

						urlVal := url
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupCommunity,
						}
						if strings.Contains(urlVal, "discord") {
							thumbImgDiscord := "fa-brands fa-discord"
							linkItem.LinkIcon = &thumbImgDiscord

							titleDiscord := "Discord"
							linkItem.LinkTitle = &titleDiscord
						} else if strings.Contains(urlVal, "facebook") {
							thumbImgFacebook := "fa-brands fa-facebook"
							linkItem.LinkIcon = &thumbImgFacebook

							titleFacebook := "Facebook"
							linkItem.LinkTitle = &titleFacebook
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}

				for _, url := range coingeckoDetailCoin.Links.OfficialForumUrl {
					//don't append blank link
					if url != "" {
						urlVal := url
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupCommunity,
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}

				if coingeckoDetailCoin.Links.BitcointalkThreadIdentifier != nil {
					urlVal := fmt.Sprintf("https://bitcointalk.org/index.php?topic=%d", *coingeckoDetailCoin.Links.BitcointalkThreadIdentifier)
					linkItem := &dao.LinkItem{
						LinkUrl:     &urlVal,
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupCommunity,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				if coingeckoDetailCoin.Links.TelegramChannelIdentifier != "" {
					urlVal := fmt.Sprintf("https://t.me/%s", coingeckoDetailCoin.Links.TelegramChannelIdentifier)
					thumbImgTelegram := "fa-brands fa-telegram"
					titleTelegram := "Telegram"
					linkItem := &dao.LinkItem{
						LinkUrl:     &urlVal,
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupCommunity,
						LinkIcon:    &thumbImgTelegram,
						LinkTitle:   &titleTelegram,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				groupSearchOnTwitter := 5
				if coinRef.CoinSymbol != nil {
					urlVal := fmt.Sprintf(`https://twitter.com/search?q=$%s`, *coinRef.CoinSymbol)
					thumbImgTwitterSearch := "fa-solid fa-magnifying-glass"
					titleTwitter := "Twitter"
					linkItem := &dao.LinkItem{
						LinkUrl:     &urlVal,
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupSearchOnTwitter,
						LinkIcon:    &thumbImgTwitterSearch,
						LinkTitle:   &titleTwitter,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				groupSourceCode := 6
				for _, url := range coingeckoDetailCoin.Links.ReposUrl.Github {
					//don't append blank link
					if url != "" {
						urlVal := url
						thumbImgGithub := "fa-brands fa-github"
						titleGithub := "Github"
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupSourceCode,
							LinkIcon:    &thumbImgGithub,
							LinkTitle:   &titleGithub,
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}

				for _, url := range coingeckoDetailCoin.Links.ReposUrl.Bitbucket {
					//don't append blank link
					if url != "" {
						urlVal := url
						thumbImgBitbucket := "fa-brands fa-bitbucket"
						titleBitbucket := "Bitbucket"
						linkItem := &dao.LinkItem{
							LinkUrl:     &urlVal,
							CoinId:      coinRef.CoinID,
							LinkGroupId: &groupSourceCode,
							LinkIcon:    &thumbImgBitbucket,
							LinkTitle:   &titleBitbucket,
						}
						_, found := existedLinkItemMap[linkItem.GetKeyMap()]
						//Append for prepare insert
						if !found {
							linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
						}
					}
				}

				groupAPIId := 7
				if coinRef.CoinID != nil {
					titleCoinId := *coinRef.CoinID
					linkItem := &dao.LinkItem{
						CoinId:      coinRef.CoinID,
						LinkGroupId: &groupAPIId,
						LinkTitle:   &titleCoinId,
					}
					_, found := existedLinkItemMap[linkItem.GetKeyMap()]
					//Append for prepare insert
					if !found {
						linkItemUrlRepo.LinkItems = append(linkItemUrlRepo.LinkItems, linkItem)
					}
				}

				linkItemUrlRepo.InsertDB()
				// insertedLinkItem := linkItemUrlRepo.InsertDB()
				// fmt.Println("2===insert success", insertedLinkItem, "linkItemContractRepo")
				//for next loop reset
				linkItemUrlRepo.LinkItems = make([]*dao.LinkItem, 0)
				//######### END: Insert link item #########

				//######### Start: Insert coin tag #########
				for _, category := range coingeckoDetailCoin.Categories {
					categoryVal := category
					coinTag := &dao.CoinTag{
						CoinId:  coinRef.CoinID,
						TagName: &categoryVal,
					}
					coinTagRepo.CoinsTags = append(coinTagRepo.CoinsTags, coinTag)
				}
				coinTagRepo.InsertDB()
				// insertTag, insertedcoinTag := coinTagRepo.InsertDB()
				// fmt.Println("3===insert success", insertTag, " insertTag; ", insertedcoinTag, "coinTag", " coin id", *coinRef.CoinID)
				//######### END: Insert coin tag #########

				// Update current coin id to old coin id for next loop
				oldCoinId = *coinRef.CoinID
			}

		}

	}

	fmt.Println("=======================================Crawl done detail==================================")

}

type CoingeckoDetailCoin struct {
	CoinId          string         //must inititalize in constructor
	Descriptions    any            `json:"description"`
	Images          any            `json:"image"`
	Links           LinkDetailCoin `json:"links"`
	Platforms       any            `json:"platforms"`
	DetailPlatforms any            `json:"detail_platforms"` //k: chainName; v:{decimal_place: any.(int); contract_address:any.(string)}
	Categories      []string       `json:"categories"`
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

func (coingeckoDTO *CoingeckoDetailCoin) Crawl() error {
	for {
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Get(fmt.Sprintf(detailCoinUrl, coingeckoDTO.CoinId))
		if err != nil { // Call API from coingecko error, or timeout, network missing.
			time.Sleep(_missRequestWait)
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
		if notExistedModel.Error != "" {
			fmt.Println("notExistedModel", notExistedModel, notExistedModel.Error)
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Not found coin from API(coin deleted) ")
			return errors.New("not found coin from API(coin deleted)")
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
				time.Sleep(_missRequestWait)
				continue
			} else { //ko loi
				return nil
			}

		} else { //Bi chan
			log.Println(log.LogLevelWarn, "detail_coins/Crawl", "Block due to time limit ")
			time.Sleep(10 * time.Second)
			continue
		}

	}
}
