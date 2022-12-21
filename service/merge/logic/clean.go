package logic

import (
	"crawler/pkg/log"
	"crawler/service/merge/model/dao"
	"fmt"
	"strings"
)

func init() {
	go AutoCrawlDataIcoHolder()
	go func() {
		//FormatSocialsDappRadar// FormatSocialsDappRadarVer2()
		// MergeDataDappradarAndRevain()
		// FormatDataDappradarWithoutRevain()
		// FormatDataRevainWithoutDappradar()
		// FormatDataProjectResearch() //Not test convert
	}()
	//Go routine insert new category, subcategory --> de subcategory query lau
}

func FormatSocialsDappRadar() {
	_src := `dappradar`
	repo := dao.ProjectRepo{}
	err := repo.SelectBySrc(_src)
	if err != nil {
		log.Println(log.LogLevelError, `service/merge/logic/clean.go/FormatSocialsDappRadar()/repo.SelectBySrc(dappradar)`, err.Error())
		return
	}
	for _, dao := range repo.Projects {
		isUpdated, err := dao.IsUpdatedSocialsByCodeAndSrc()
		if err != nil {
			log.Println(log.LogLevelError, `service/merge/logic/clean.go/FormatSocialsDappRadar()/dao.IsUpdatedSocialsByCodeAndSrc()`, err.Error())
			continue
		}
		if isUpdated {
			continue
		}
		idx := 0
		newSocials := make(map[string]any)
		for socialUrl, socialImg := range dao.Socials {
			data := struct {
				key          string
				Tittle       string `json:"tittle"`
				Url          string `json:"url"`
				SvgHtmlImage string `json:"svgHtmlImage"`
			}{}
			//Twitter
			if strings.Contains(socialUrl, `twitter.com`) {
				data.key = `twitter`
				data.Tittle = `Twitter`
			} else
			//Facebook
			if strings.Contains(socialUrl, `facebook.com`) {
				data.key = `facebook`
				data.Tittle = `Facebook`
			} else
			//Discord
			if strings.Contains(socialUrl, `discord.com`) || strings.Contains(socialUrl, `discord.gg`) {
				data.key = `discord`
				data.Tittle = `Discord`
			} else
			//Reddit
			if strings.Contains(socialUrl, `reddit.com`) {
				data.key = `reddit`
				data.Tittle = `Reddit`
			} else
			//Telegram
			if strings.Contains(socialUrl, `t.me`) || strings.Contains(socialUrl, `telegram.com`) {
				data.key = `telegram`
				data.Tittle = `Telegram`
			} else
			//Instagram
			if strings.Contains(socialUrl, `instagram.com`) {
				data.key = `instagram`
				data.Tittle = `Instagram`
			} else
			//Medium
			if strings.Contains(socialUrl, `medium.com`) {
				data.key = `medium`
				data.Tittle = `Medium`
			} else
			//Github
			if strings.Contains(socialUrl, `github.com`) {
				data.key = `github`
				data.Tittle = `Github`
			} else
			//Instagram
			if strings.Contains(socialUrl, `instagram.com`) {
				data.key = `instagram`
				data.Tittle = `Instagram`
			} else
			//Youtube
			if strings.Contains(socialUrl, `youtube.com`) {
				data.key = `youtube`
				data.Tittle = `Youtube`
			} else
			//Other
			{
				if idx == 0 {
					data.key = "website"
				} else
				// Greater than 1
				{
					data.key = fmt.Sprintf("website%d", idx)
				}
				idx++
				data.Tittle = `Website`
			}
			data.Url = socialUrl
			data.SvgHtmlImage = socialImg.(string)
			newSocials[data.key] = data
		}
		dao.Socials = newSocials
		//code and src from get from db already
		err = dao.UpdateSocialsAndExtraDataByCodeAndSrc()
		if err != nil {
			log.Println(log.LogLevelError, `service/merge/logic/clean.go/FormatSocialsDappRadar()/dao.UpdateSocialsByCode() at project code: `+dao.ProjectCode, err.Error())
			continue
		}
		err = dao.InsertUpdatedSocial()
		//###############Fatal########################
		if err != nil {
			log.Println(log.LogLevelFatal, `service/merge/logic/clean.go/FormatSocialsDappRadar()/dao.InsertUpdatedSocial() at project code: `+dao.ProjectCode, err.Error())
		}

	}

}

func FormatSocialsDappRadarVer2() {
	_src := `dappradar`
	repo := dao.ProjectRepo{}
	err := repo.SelectBySrc(_src)
	if err != nil {
		log.Println(log.LogLevelError, `service/merge/logic/clean.go/FormatSocialsDappRadar()/repo.SelectBySrc(dappradar)`, err.Error())
		return
	}
	for _, dao := range repo.Projects {
		isUpdated, err := dao.IsUpdatedSocialsByCodeAndSrc()
		if err != nil {
			log.Println(log.LogLevelError, `service/merge/logic/clean.go/FormatSocialsDappRadar()/dao.IsUpdatedSocialsByCodeAndSrc()`, err.Error())
			continue
		}
		if isUpdated {
			continue
		}
		newSocials := make(map[string]any)
		newSocials[`social`] = make([]string, 0)
		for dataKey, data := range dao.Socials {
			//dataKey:= twitter, facebook, discord, reddit, telegram, instagram, medium, github, youtube, website, websiteN

			if dataKey == `github` {
				if dao.ExtraData == nil {
					dao.ExtraData = make(map[string]any)
				}
				dao.ExtraData[`sourceCode`] = data.(map[string]any)[`url`]

			} else
			//discord or url
			if strings.Contains(dataKey, `website`) {
				url, foundUrl := data.(map[string]any)[`url`]
				if foundUrl {
					if strings.Contains(url.(string), `discord.gg`) {
						url, foundUrl := data.(map[string]any)[`url`]
						if foundUrl {
							newSocials[`social`] = append(newSocials[`social`].([]string), url.(string))
						}
					} else {
						if dao.ExtraData == nil {
							dao.ExtraData = make(map[string]any)
						}
						dao.ExtraData[`website`] = url
					}
				}
			} else
			//Social
			{
				url, foundUrl := data.(map[string]any)[`url`]
				if foundUrl {
					newSocials[`social`] = append(newSocials[`social`].([]string), url.(string))
				}
			}

		}
		dao.Socials = newSocials
		//code and src from get from db already
		err = dao.UpdateSocialsAndExtraDataByCodeAndSrc()
		if err != nil {
			log.Println(log.LogLevelError, `service/merge/logic/clean.go/FormatSocialsDappRadar()/dao.UpdateSocialsByCode() at project code: `+dao.ProjectCode, err.Error())
			continue
		}
		err = dao.InsertUpdatedSocial()
		//###############Fatal########################
		if err != nil {
			log.Println(log.LogLevelFatal, `service/merge/logic/clean.go/FormatSocialsDappRadar()/dao.InsertUpdatedSocial() at project code: `+dao.ProjectCode, err.Error())
		}

	}

}

func MergeDataDappradarAndRevain() {
	projectRepo := &dao.ProjectRepo{}
	err := projectRepo.SelectMergeCommonDataFromDappAndRevain()
	if err != nil {
		log.Println(log.LogLevelError, `service/merge/logic/clean.go/MergeData()/projectRepo.SelectMergeCommonDataFromDappAndRevain()`, err.Error())
	}
	sources := map[string]bool{
		`dappradar`: true,
		`revain`:    true,
	}
	productRepo := &dao.ProductRepo{}
	projectRepo.ConvertToProductRepo(productRepo, sources)
	productRepo.InsertDB(sources)
	fmt.Println(`done insert MergeDataDappradarAndRevain`)
}

func FormatDataDappradarWithoutRevain() {
	projectRepo := &dao.ProjectRepo{}
	err := projectRepo.SelectDappradarWithoutRevain()
	if err != nil {
		log.Println(log.LogLevelError, `service/merge/logic/clean.go/MergeData()/projectRepo.SelectDappradarWithoutRevain()`, err.Error())
	}
	productRepo := &dao.ProductRepo{}
	sources := map[string]bool{
		`dappradar`: true,
	}
	projectRepo.ConvertToProductRepo(productRepo, sources)
	productRepo.FormatCategoryDappRadar()
	productRepo.InsertDB(sources)
	fmt.Println(`done insert FormatDataDappradarWithoutRevain`)
}

func FormatDataRevainWithoutDappradar() {
	projectRepo := &dao.ProjectRepo{}
	err := projectRepo.SelectRevainWithoutDappradar()
	if err != nil {
		log.Println(log.LogLevelError, `service/merge/logic/clean.go/MergeData()/projectRepo.FormatDataRevainWithoutDappradar()`, err.Error())
	}
	productRepo := &dao.ProductRepo{}
	sources := map[string]bool{
		`revain`: true,
	}
	projectRepo.ConvertToProductRepo(productRepo, sources)
	productRepo.FormatProductNameRevain()
	productRepo.InsertDB(sources)
	fmt.Println(`done insert FormatDataRevainWithoutDappradar`)
}

func FormatDataProjectResearch() {
	projectResearchRepo := &dao.ProjectResearchRepo{}
	projectResearchRepo.SelectAll()

	productRepo := &dao.ProductRepo{}
	projectResearchRepo.ConvertToProductRepo(productRepo)

	sources := map[string]bool{
		`research`: true,
	}

	productRepo.InsertDB(sources)
}
