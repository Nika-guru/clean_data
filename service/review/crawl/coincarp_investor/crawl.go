package crawl_coincarp_investor

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Crawl() {
	for {
		CrawlFundRaisings()
		time.Sleep(1 * time.Hour)
	}
}

func CrawlFundRaisings() {
	maxGoroutines := 1
	guard := make(chan struct{}, maxGoroutines)
	for idx := 0; idx < 5; idx++ {
		guard <- struct{}{} //buffered channel, full capacity, wait here --> limit go routine

		go func(idx int) {
			fmt.Println(`====page`, (idx + 1))
			CrawlFundRaisingsByPageIndex(idx)
			<-guard
		}(idx)

	}
}

func CrawlFundRaisingsByPageIndex(idx int) {
	url := fmt.Sprintf(constant.URL_GET_FUNDRAISINGS_COINCARP, idx*constant.LIMIT_REC_FUND_RAISINGS_COINCARP, constant.LIMIT_REC_FUND_RAISINGS_COINCARP)

	dtoResp := FundCoincarpResp{}
	err := dtoResp.CrawlFund(url)
	if err != nil {
		log.Println(log.LogLevelError, "service/review/crawl/coincarp_investor/crawl.go/CrawlFundRaisings/dtoResp.CrawlFund(url)", err.Error())
		time.Sleep(10 * time.Second)
		return
	}

	//Get list investors of current fund raising for each project
	for _, fundItem := range dtoResp.FundData.FundItems {
	HTMLFundMoreData:
		fundItem.SrcFund += fmt.Sprintf("%s\n", url)

		//Crawl fund info(description, announcementurl)
		url := constant.BASE_URL_COINCARP + constant.ENDPOINT_FUND_RAISING_COINCARP + `/` + fundItem.EndpointFund
		dom := utils.GetHtmlDomByUrl(url)
		hasData := CrawlDescriptionAndAnnouncementByHtmlDom(dom, &fundItem)
		if !hasData {
			goto HTMLFundMoreData
		}
		fundItem.SrcFund += fmt.Sprintf("%s\n", url)

		//Crawl fund info(investor info)
		fundRepo := dao.FuncRaisingRepo{}
		fundRepo.FuncRaisingList = CrawlInvestorsByFundRasing(fundItem)
		fundRepo.InsertDB()

		//#######################################################################################################################
		//Crawl invertors info
		for _, daoFundRaising := range fundRepo.FuncRaisingList {
			//don't query API get investor not exist
			if daoFundRaising.InvestorCode == `NULL` {
				continue
			}
			investorRepo := dao.InvestorRepo{}
			investor := &dao.Investor{}
			investor.InvestorCode = daoFundRaising.InvestorCode

			//check investors code exist ib here
			isExist, err := investor.IsExist()
			if err != nil {
				log.Println(log.LogLevelError, "service/review/crawl/coincarp_investor/crawl.go/CrawlFundRaisings/investor.IsExist()/investor code: "+investor.InvestorCode, err.Error())
			} else {
				if !isExist {
					url = constant.BASE_URL_COINCARP + constant.ENDPOINT_INVESTOR_COINCARP + `/` + investor.InvestorCode
				HTMLInvestorDetail:
					dom := utils.GetHtmlDomByUrl(url)
					investor.Src = url
					hasData := CrawlInvestorByHtmlDom(dom, investor)
					if !hasData {
						goto HTMLInvestorDetail
					}

					investorRepo.Investors = append(investorRepo.Investors, *investor)
				}
				investorRepo.InsertDB()
			}

		}
		//#######################################################################################################################

		continue
		//#######################################################################################################################
		//Crawl project/coin info
		url = fundItem.ProjectCode
		daoProject := &dao.Project{}
		daoCoin := &dao.Coin{}
	HTMLProjectDetail:
		dom = utils.GetHtmlDomByUrl(url)
		hasData, isCoin := CrawlProjectByHtmlDom(dom, daoProject, daoCoin)
		if !hasData {
			goto HTMLProjectDetail
		}
		if isCoin {
			daoCoin.CoinId = fundItem.ProjectCode
			daoCoin.Src = url
			err := daoCoin.InsertDB()
			if err != nil {
			}
		} else {
			daoProject.ProjectCode = fundItem.ProjectCode
			daoProject.Src = url
			err := daoProject.InsertDB()
			if err != nil {
			}
		}
		//#######################################################################################################################

	}
}

func CrawlInvestorsByFundRasing(fundItem FundItem) []dao.FuncRaising {
	daos := make([]dao.FuncRaising, 0)

	hasInvestor := false
	for idx := 0; ; idx++ {
		fundCode := fundItem.EndpointFund //projectCode + '-' + fundStageCode
		url := fmt.Sprintf(constant.URL_GET_INVESTOR_BY_FUND_RAISING_COINCARP, fundItem.ProjectCode, fundCode, (idx + 1), constant.LIMIT_INVESTOR_BY_FUND_RAISING_COINCARP)

		dtoResp := InvestorCoincarpResp{}
		finish, err := dtoResp.CrawlInvestor(url)
		if err != nil {
			//Thuong la loi goi tin http, khong unmarshal duoc --> goi lai page --> log thi nhieu qua
			// log.Println(log.LogLevelWarn, "service/review/crawl/coincarp_investor/crawl.go/CrawlInvestorsByFundRasing/dtoResp.CrawlInvestor(url)", err.Error())
			idx--
			time.Sleep(5 * time.Second)
			continue
		}

		//Crawl finish all investors
		if finish {
			break
		}
		for _, investorItem := range dtoResp.InvestorData.InvestorItems {
			hasInvestor = true
			daos = append(daos, dao.FuncRaising{
				ProjectCode:     fundItem.ProjectCode,
				ProjectName:     fundItem.ProjectName,
				ProjectLogo:     constant.BASE_URL_IMAGE_SERVER_COINCARP + fundItem.ProjectLogo,
				InvestorCode:    investorItem.InvestorCode,
				InvestorName:    investorItem.InvestorName,
				InvestorLogo:    investorItem.Logo,
				FundStageCode:   fundItem.FundStageCode,
				FundStageName:   fundItem.FundStageName,
				FundAmount:      fundItem.FundAmount,
				FundDate:        utils.TimestampToString(time.Unix(int64(fundItem.FundDate), 0).UTC()),
				Description:     fundItem.Description,
				AnnouncementUrl: fundItem.AnnouncementUrl,
				Valulation:      fundItem.Valulation,
				SrcInvestor:     url,
				SrcFund:         fundItem.SrcFund,
			})
		}
	}

	//ink-games don't have any data of investors for series-a, series-b
	if !hasInvestor {
		daos = append(daos, dao.FuncRaising{
			ProjectCode:     fundItem.ProjectCode,
			ProjectName:     fundItem.ProjectName,
			ProjectLogo:     constant.BASE_URL_IMAGE_SERVER_COINCARP + fundItem.ProjectLogo,
			InvestorCode:    `NULL`,
			InvestorName:    `NULL`,
			InvestorLogo:    `NULL`,
			FundStageCode:   fundItem.FundStageCode,
			FundStageName:   fundItem.FundStageName,
			FundAmount:      fundItem.FundAmount,
			FundDate:        utils.TimestampToString(time.Unix(int64(fundItem.FundDate), 0).UTC()),
			Description:     fundItem.Description,
			AnnouncementUrl: fundItem.AnnouncementUrl,
			Valulation:      fundItem.Valulation,
			SrcInvestor:     `NULL`,
			SrcFund:         fundItem.SrcFund,
		})
	}
	return daos
}

func CrawlDescriptionAndAnnouncementByHtmlDom(dom *goquery.Document, fundItem *FundItem) (hasData bool) {
	description := ``
	domKey := `h3` + utils.ConvertClassesFormatFromBrowserToGoQuery(`font-size-18`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		description += s.Text()
		description += "\n"
	})

	domKey = `div#projectInfo`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `p`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			description += s.Text()
			description += "\n"
		})

	})
	fundItem.Description = description

	announcementUrl := ``
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`font-size-12 text-grey mt-3`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		domKey = `a`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			url, foundUrl := s.Attr(`href`)
			if foundUrl {
				announcementUrl = url
			}

		})
	})
	fundItem.AnnouncementUrl = announcementUrl

	return (description != ``) //&& announcementUrl != `` https://www.coincarp.com/fundraising/recnlembdnnfyk/
}

func CrawlInvestorByHtmlDom(dom *goquery.Document, investor *dao.Investor) (hasData bool) {
	hasData = true
	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`funding-detail-top d-flex align-items-center`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `img`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			img, foundImg := s.Attr(`src`)
			if foundImg {
				investor.InvestorImage = img
			} else {
				hasData = false
			}
		})

		hasData = false
		domKey = `h1`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			investor.InvestorName = s.Text()
			hasData = true
		})

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`social-list d-flex`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			if investor.Socials == nil {
				investor.Socials = make(map[string]any, 0)
			}

			domKey = `a`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				url, foundUrl := s.Attr("href")
				if foundUrl {
					s.Each(func(i int, s *goquery.Selection) {
						childTag, _ := s.Html()
						socialData := struct {
							key         string
							HtmlIconTag string `json:"htmlIconTag"`
							Title       string `json:"title"`
							Url         string `json:"url"`
						}{}
						//siteLink
						if strings.Contains(childTag, `<i class="iconfont icon-link font-size-14 font-weight-normal">`) {
							socialData.key = `siteLink`
							socialData.Title = `Official Website`
						} else
						//twitter
						if strings.Contains(childTag, `<i class="iconfont icon-twitter font-size-14 font-weight-normal">`) {
							socialData.key = `twitter`
							socialData.Title = `Twitter`
						} else
						//linkedin
						if strings.Contains(childTag, `<i class="iconfont icon-in font-size-14 font-weight-normal">`) {
							socialData.key = `linkedin`
							socialData.Title = `Linkedin`
						} else
						//email
						if strings.Contains(childTag, `<svg`) {
							socialData.key = `email`
							socialData.Title = `Email`
						}
						socialData.Url = url
						socialData.HtmlIconTag = childTag
						investor.Socials[socialData.key] = socialData
					})
				}
			})
		})
	})

	domKey = `span` + utils.ConvertClassesFormatFromBrowserToGoQuery(`catagory-tag d-inline-block font-size-12`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		investor.CategoryName = s.Text()
	})

	domKey = `p` + utils.ConvertClassesFormatFromBrowserToGoQuery(`font-weight-bold text-dark`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		//Location
		if i == 0 {
			investor.Location = s.Text()
		} else
		//Year Founded
		if i == 1 {
			yearFounded, err := strconv.Atoi(s.Text())
			//No data year
			if err != nil {
				investor.YearFounded = 0
			} else {
				investor.YearFounded = yearFounded
			}
		}
	})

	domKey = `div#projectInfo`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		investor.Description, _ = s.Html()
	})

	return hasData
}

func CrawlProjectByHtmlDom(dom *goquery.Document, daoProject *dao.Project, daoCoin *dao.Coin) (hasData bool, isCoin bool) {
	hasData = true

	projectName := `` //
	projectLogo := `` //
	yearLaunched := 0 //
	location := ``    //
	category := ``    //
	subcategory := `` //
	description := `` //
	founders := make([]any, 0)
	socials := make(map[string]any, 0) //

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`funding-detail-top d-flex align-items-center`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `img`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			img, foundImg := s.Attr(`src`)
			if foundImg {
				projectLogo = img
			} else {
				hasData = false
			}
		})

		hasData = false
		domKey = `h1`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			projectName = s.Text()
			hasData = true
		})

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`social-list d-flex`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			if socials == nil {
				socials = make(map[string]any, 0)
			}

			domKey = `a`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				url, foundUrl := s.Attr("href")
				if foundUrl {
					s.Each(func(i int, s *goquery.Selection) {
						childTag, _ := s.Html()
						socialData := struct {
							key         string
							HtmlIconTag string `json:"htmlIconTag"`
							Title       string `json:"title"`
							Url         string `json:"url"`
						}{}
						//siteLink
						if strings.Contains(childTag, `<i class="iconfont icon-link font-size-14 font-weight-normal">`) {
							socialData.key = `siteLink`
							socialData.Title = `Official Website`
						} else
						//twitter
						if strings.Contains(childTag, `<i class="iconfont icon-twitter font-size-14 font-weight-normal">`) {
							socialData.key = `twitter`
							socialData.Title = `Twitter`
						} else
						//linkedin
						if strings.Contains(childTag, `<i class="iconfont icon-in font-size-14 font-weight-normal">`) {
							socialData.key = `linkedin`
							socialData.Title = `Linkedin`
						} else
						//email
						if strings.Contains(childTag, `<svg`) {
							socialData.key = `email`
							socialData.Title = `Email`
						}
						socialData.Url = url
						socialData.HtmlIconTag = childTag
						socials[socialData.key] = socialData
					})
				}
			})
		})
	})

	domKey = `span` + utils.ConvertClassesFormatFromBrowserToGoQuery(`catagory-tag d-inline-block font-size-12`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			category = s.Text()
		} else if i == 1 {
			subcategory = s.Text()
		}
	})

	domKey = `p` + utils.ConvertClassesFormatFromBrowserToGoQuery(`font-weight-bold text-dark`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		//Location
		if i == 0 {
			location = s.Text()
		} else
		//Year Founded
		if i == 1 {
			yearFounded, err := strconv.Atoi(s.Text())
			//No data year
			if err != nil {
				yearLaunched = 0
			} else {
				yearLaunched = yearFounded
			}
		}
	})

	domKey = `div#projectInfo`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		description, _ = s.Html()
	})

	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`founders-list d-flex flex-wrap`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		//append in this for
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`item text-center border rounded`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			FouderNameVal := ``
			domKey = `p`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				FouderNameVal = s.Text()
			})

			FouderImgVal := ``
			domKey = `img`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				img, foundImg := s.Attr(`src`)
				if foundImg {
					FouderImgVal = img
				}
			})

			FounderSocialsVal := make(map[string]any, 0)
			//
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`social-list d-flex justify-content-center`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					url, foundUrl := s.Attr("href")
					if foundUrl {
						s.Each(func(i int, s *goquery.Selection) {
							childTag, _ := s.Html()
							socialData := struct {
								key         string
								HtmlIconTag string `json:"htmlIconTag"`
								Title       string `json:"title"`
								Url         string `json:"url"`
							}{}
							//siteLink
							if strings.Contains(childTag, `<i class="iconfont icon-link font-size-14 font-weight-normal">`) {
								socialData.key = `siteLink`
								socialData.Title = `Official Website`
							} else
							//twitter
							if strings.Contains(childTag, `<i class="iconfont icon-twitter font-size-14 font-weight-normal">`) {
								socialData.key = `twitter`
								socialData.Title = `Twitter`
							} else
							//linkedin
							if strings.Contains(childTag, `<i class="iconfont icon-in font-size-14 font-weight-normal">`) {
								socialData.key = `linkedin`
								socialData.Title = `Linkedin`
							} else
							//email
							if strings.Contains(childTag, `<svg`) {
								socialData.key = `email`
								socialData.Title = `Email`
							}
							socialData.Url = url
							socialData.HtmlIconTag = childTag
							FounderSocialsVal[socialData.key] = socialData
						})
					}
				})
			})

			founders = append(founders, struct {
				FouderName     string
				FouderImg      string
				FounderSocials map[string]any
			}{
				FouderName:     FouderNameVal,
				FouderImg:      FouderImgVal,
				FounderSocials: FounderSocialsVal,
			})
		})
	})

	isCoinToken := false
	domKey = `h2`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		if s.Text() == `Coin or Token` {
			isCoinToken = true
		}
	})
	if isCoinToken {
		daoCoin.Name = projectName
		daoCoin.Tag = strings.TrimSpace(category)
		if daoCoin.Tag != `` && strings.TrimSpace(subcategory) != `` {
			daoCoin.Tag += `,`
		}
		daoCoin.Tag = strings.TrimSpace(subcategory)
		daoCoin.Image = projectLogo
		if daoCoin.Detail == nil {
			daoCoin.Detail = make(map[string]any, 0)
		}
		daoCoin.Detail[`description`] = description
		daoCoin.Detail[`socials`] = socials
		daoCoin.Detail[`founders`] = founders
		daoCoin.Detail[`location`] = location
		daoCoin.Detail[`yearLaunched`] = yearLaunched
	} else {

	}
	return false, false
}
