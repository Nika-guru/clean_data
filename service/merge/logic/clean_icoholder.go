package logic

import (
	"crawler/pkg/log"
	"crawler/pkg/utils"
	"crawler/service/merge/model/dao"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	_baseUrl                     = `https://icoholder.com`
	_endpointOngoingListByPage   = `/en/icos/ongoing?isort=r.general&idirection=desc&page=`
	_endpointUpcommingListByPage = `/en/icos/upcoming?isort=r.general&idirection=desc&page=`
	_endpointHome                = `/en/cryptos?sort=q.market_cap&direction=desc&page=`
)

func AutoCrawlDataIcoHolder() {
	// go CrawlAllList(_endpointHome)
	for {
		CrawlAllList(_endpointOngoingListByPage)
		CrawlAllList(_endpointUpcommingListByPage)
		time.Sleep(12 * time.Hour)
	}
}

func CrawlAllList(endpointList string) {
	maxGoroutines := 2
	guard := make(chan struct{}, maxGoroutines)
	for idx := 0; ; idx++ {
		guard <- struct{}{} //buffered channel, full capacity, wait here --> limit go routine

		go func(endpointList string, idx int) {
			endpoints := CrawlOngoingListByPagination(endpointList, idx)

			//last page have data
			if len(endpoints) == 0 {
				return
			}

			for _, endpoint := range endpoints {
				urlDetail := fmt.Sprintf("%s%s", _baseUrl, endpoint)
				product, members := CrawlDetailByUrl(urlDetail, true)
				chainList := dao.ChainList{}
				chainList.ChainName = product.ChainName
				chainList.SelectChainIdByChainName()
				if chainList.ChainId == `` {
					product.ChainId = `NULL`
				} else {
					product.ChainId = chainList.ChainId
				}
				productRepo := dao.ProductRepo{}
				productRepo.Products = append(productRepo.Products, &product)
				productRepo.InsertDB(map[string]bool{
					`icoholder`: true,
				})
				productId := productRepo.Products[0].Id

				memberRepo := dao.MemberRepo{}
				//update product id
				for _, member := range members {
					member.ProductId = uint64(productId)
					memberRepo.Members = append(memberRepo.Members, &member)
				}
				memberRepo.InsertDB()
			}

			<-guard
		}(endpointList, idx)
	}
}

func CrawlByEndpointListAndIndex() {

}

func CrawlOngoingListByPagination(endpointList string, idx int) (endpoints []string) {
	url := fmt.Sprintf("%s%s%d", _baseUrl, endpointList, (idx + 1))
call:
	dom := utils.GetHtmlDomJsRenderByUrl(url)

	if dom == nil {
		log.Println(log.LogLevelWarn, `service/merge/logic/clean_icoholder.go/CrawlListByPagination/utils.GetHtmlDomJsRenderByUrl(url)`, `dom nil`)
		time.Sleep(5 * time.Second)
		goto call
	}
	endpoints = make([]string, 0)

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`list-ico-row`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		exceptEndpoint := make(map[string]bool, 0)
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-list-row ico-stared`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `h3`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					url, foundUrl := s.Attr(`href`)
					if foundUrl {
						exceptEndpoint[url] = true
					}
				})

			})

		})

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-list-row`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			domKey = `h3`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					url, foundUrl := s.Attr(`href`)
					if foundUrl {
						_, found := exceptEndpoint[url]
						if !found {
							endpoints = append(endpoints, url)
						}
					}
				})

			})
		})

	})

	return endpoints
}

func CrawlDetailByUrl(url string, isICO bool) (product dao.Product, members []dao.Member) {
	product = dao.Product{}
	if product.Detail == nil {
		product.Detail = make(map[string]any)
	}
	product.Category = `Crypto Projects`
	if isICO {
		product.Type = `ico`
	}

	//For member struct
	productName := ``
	productSymbol := ``
	members = make([]dao.Member, 0)
call:
	fmt.Println(url)
	dom := utils.GetHtmlDomJsRenderByUrl(url)

	if dom == nil {
		log.Println(log.LogLevelWarn, `service/merge/logic/clean_icoholder.go/CrawlListByPagination/utils.GetHtmlDomJsRenderByUrl(url)`, `dom nil`)
		time.Sleep(5 * time.Second)
		goto call
	}

	parts := strings.Split(url, `/`)
	if len(parts) >= 2 {
		productCode := parts[len(parts)-1]
		product.Detail[`productCodeIcoholder`] = productCode

		//Convert icoholder coinid to coingecko coinid
		// parts = strings.Split(productCode, `-`)
		// productCode = ``
		// for idx, part := range parts {
		// 	if idx < len(parts)-1 {
		// 		productCode += part
		// 	}
		// 	if idx < len(parts)-2 {
		// 		productCode += `-`
		// 	}
		// }
		// fmt.Println(`===`, parts, productCode)
	}

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`logo-name logo-name-view col-md-12`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		//Image
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`col-md-1 col-xs-3`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey := `img`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				img, fountImg := s.Attr(`data-src`)
				if fountImg {
					product.Image = strings.TrimSpace(img)
				}
			})

		})

		//Name
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-titles-in-view`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			domKey = `h1`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				product.Name = strings.TrimSpace(s.Text())
				productName = product.Name
			})
		})

		//Description
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`description-value`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			product.Description = strings.TrimSpace(s.Text())
		})

	})

	//Chain name, Subcategories
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info ico-more-info--second`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info__row`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			isPlatform := false
			isCategories := false
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info__subtitle`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				if s.Text() == `Platform` {
					isPlatform = true
				} else if s.Text() == `Categories` {
					isCategories = true
				}
			})

			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info__opacity`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				if isPlatform {
					product.ChainName = strings.ReplaceAll(s.Text(), ` `, ``)
					product.ChainName = strings.ToLower(product.ChainName)
					if !isICO {
						if product.ChainName == `` {
							product.Type = `coin`
						} else {
							product.Type = `token`
						}
					}
				} else if isCategories {
					product.Subcategory = strings.TrimSpace(s.Text())
				}
			})

		})

	})

	//Symbol
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`cg-base cg-ticker`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		pair := s.Text()
		parts := strings.Split(pair, `/`)
		if len(parts) == 2 {
			product.Symbol = strings.TrimSpace(parts[0])
		}
	})

	//Symbol (sub) + Max supply
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info ico-more-info--first`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info__row`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			isSymbol := false
			isTotalSupply := false
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info__subtitle`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				if s.Text() == `Ticker` {
					isSymbol = true
				} else if s.Text() == `Total supply` {
					isTotalSupply = true
				}
			})

			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`ico-more-info__opacity`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				if isSymbol {
					if product.Symbol == `` {
						product.Symbol = strings.TrimSpace(strings.ToUpper(s.Text()))
						productSymbol = product.Symbol
					}
				} else if isTotalSupply {
					// maxSupplyFloat, _ := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(s.Text()), `,`, ``), 64) --> co nhung truong hop 200M MAI, ...
					product.Detail[`maxSupply`] = strings.TrimSpace(s.Text())
				}
			})

		})

	})

	//About(+description)
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`col-md-12 col-sm-12 about-value`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		if product.Description != `` {
			product.Description += "\n"
		}
		product.Description += strings.TrimSpace(s.Text())
	})

	//Social + Website
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`links-right`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		//Social
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`project-links`)
		dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
			domKey = `a`
			socials := make([]string, 0)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				url, foundUrl := s.Attr(`href`)
				if foundUrl {
					socials = append(socials, url)
				}
			})
			product.Detail[`social`] = socials
		})

		//Website
		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`text-align-center`)
		dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
			domKey = `a`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				url, foundUrl := s.Attr(`href`)
				if foundUrl {
					product.Detail[`website`] = url
				}
			})
		})
	})

	//Member(team)
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`col-md-12 col-sm-12 members-row`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`col-lg-4 col-sm-6 text-center mb-4`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			member := dao.Member{}
			if member.Detail == nil {
				member.Detail = make(map[string]any)
			}
			member.Detail[`positionGroup`] = `Team`
			member.Detail[`src`] = url
			member.ProductName = productName
			member.ProductSymbol = productSymbol

			//Member Image
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-icon`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				//Sample: background: url('https://icoholder.com/media/cache/member_thumb/assets/app/img/need-verify.png');
				style, foundStyle := s.Attr(`style`)
				if foundStyle {
					parts := strings.Split(style, "'")
					if len(parts) == 3 {
						member.Detail[`memberImage`] = strings.TrimSpace(parts[1])
					}
				}
			})

			//Member Linkedin
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-links`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					link, foundLink := s.Attr(`href`)
					if foundLink {
						member.Detail[`memberLinkedin`] = strings.TrimSpace(link)
					}
				})
			})

			//Member Name
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-title`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `span`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					member.MemberName = strings.TrimSpace(s.Text())
				})

			})

			//Member Position
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-position`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				member.Detail[`memberPosition`] = strings.TrimSpace(s.Text())
			})

			//IsVerified
			isVerified := false
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`verified-line verified`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				isVerified = true
			})
			member.Detail[`isVerified`] = isVerified

			members = append(members, member)
		})

	})

	//Member(advisors)
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`row advisers-row detail-padding20`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`col-lg-4 col-sm-6 text-center mb-4`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			member := dao.Member{}
			if member.Detail == nil {
				member.Detail = make(map[string]any)
			}
			member.Detail[`positionGroup`] = `Advisors`
			member.Detail[`src`] = url
			member.ProductName = productName
			member.ProductSymbol = productSymbol

			//Member Image
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-icon`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				//Sample: background: url('https://icoholder.com/media/cache/member_thumb/assets/app/img/need-verify.png');
				style, foundStyle := s.Attr(`style`)
				if foundStyle {
					parts := strings.Split(style, "'")
					if len(parts) == 3 {
						member.Detail[`memberImage`] = strings.TrimSpace(parts[1])
					}
				}
			})

			//Member Linkedin
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-links`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					link, foundLink := s.Attr(`href`)
					if foundLink {
						member.Detail[`memberLinkedin`] = strings.TrimSpace(link)
					}
				})
			})

			//Member Name
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-title`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `span`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					member.MemberName = strings.TrimSpace(s.Text())
				})

			})

			//Member Position
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-position`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				member.Detail[`memberPosition`] = strings.TrimSpace(s.Text())
			})

			//IsVerified
			isVerified := false
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`verified-line verified`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				isVerified = true
			})
			member.Detail[`isVerified`] = isVerified

			members = append(members, member)
		})

	})

	//Member(former)
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`row past-row detail-padding20`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`col-lg-4 col-sm-6 text-center mb-4`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			member := dao.Member{}
			if member.Detail == nil {
				member.Detail = make(map[string]any)
			}
			member.Detail[`positionGroup`] = `Former Members`
			member.Detail[`src`] = url
			member.ProductName = productName
			member.ProductSymbol = productSymbol

			//Member Image
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-icon`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				//Sample: background: url('https://icoholder.com/media/cache/member_thumb/assets/app/img/need-verify.png');
				style, foundStyle := s.Attr(`style`)
				if foundStyle {
					parts := strings.Split(style, "'")
					if len(parts) == 3 {
						member.Detail[`memberImage`] = strings.TrimSpace(parts[1])
					}
				}
			})

			//Member Linkedin
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-links`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					link, foundLink := s.Attr(`href`)
					if foundLink {
						member.Detail[`memberLinkedin`] = strings.TrimSpace(link)
					}
				})
			})

			//Member Name
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-title`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `span`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					member.MemberName = strings.TrimSpace(s.Text())
				})

			})

			//Member Position
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`member-position`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				member.Detail[`memberPosition`] = strings.TrimSpace(s.Text())
			})

			//IsVerified
			isVerified := false
			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`verified-line verified`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				isVerified = true
			})
			member.Detail[`isVerified`] = isVerified

			members = append(members, member)
		})

	})

	return product, members
}
