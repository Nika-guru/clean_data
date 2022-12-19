package dao

import (
	"base/pkg/db"
	"encoding/json"
	"fmt"
	"strings"
)

type ProjectResearchRepo struct {
	ProjectResearchList []*ProjectResearch
}

type ProjectResearch struct {
	Id             string //Auto incremental --> for delete after convert
	Type           string
	Address        string
	ChainId        string
	Symbol         string
	Name           string
	Category       string
	Subcategory    string
	TotalSupply    string
	MaxSupply      string
	Marketcap      string
	VolumneTrading string
	Image          string
	Decimals       string
	Detail         map[string]any
	//======skip======
	//CreatedDate string
	//UserId int
	//UserRole int
	//Username string
}

func (repo *ProjectResearchRepo) SelectAll() error {
	query := `
	SELECT 
		id, "type", address, 
		chainid, symbol, "name",
		category, subcategory, totalsupply,
		maxsupply, marketcap, volumetrading,
		image, decimals, detail
	FROM project_research;
	`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		dao := &ProjectResearch{}

		detailJSON := []byte{}

		err := rows.Scan(&dao.Id, &dao.Type, &dao.Address,
			&dao.ChainId, &dao.Symbol, &dao.Name,
			&dao.Category, &dao.Subcategory, &dao.TotalSupply,
			&dao.MaxSupply, &dao.Marketcap, &dao.VolumneTrading,
			&dao.Image, &dao.Decimals, &detailJSON)
		if err != nil {
			return err
		}
		err = json.Unmarshal(detailJSON, &dao.Detail)
		if err != nil {
			return err
		}

		repo.ProjectResearchList = append(repo.ProjectResearchList, dao)
	}

	return nil
}

func (repo *ProjectResearchRepo) ConvertToProjectAndCoinRepo(projectRepo *ProjectRepo, coinRepo *CoinRepo) {
	deletedIds := make([]string, 0)
	for _, projectResearch := range repo.ProjectResearchList {
		switch strings.TrimSpace(projectResearch.Type) {
		//Coin/ token/ ico
		case `coin`:
			fallthrough
		case `token`:
			fallthrough
		case `ico`:
			dao := &Coin{}
			coinRepo.Coins = append(coinRepo.Coins, dao)

			deletedIds = append(deletedIds, projectResearch.Id)
			//dApp
		case `project`:
			dao := &Project{}
			dao.ProjectCode = ``
			dao.ProjectName = projectResearch.Name
			dao.ProjectCategory = projectResearch.Category
			dao.ProjectSubcategory = projectResearch.Subcategory

			//avoid nil pointer below
			if projectResearch.Detail == nil {
				projectResearch.Detail = make(map[string]any, 0)
			}

			community, foundCommunity := projectResearch.Detail[`community`]
			if foundCommunity && len(community.(map[string]any)) > 0 {
				dao.Socials = make(map[string]any, 0)
				dao.Socials[`social`] = make([]string, 0)
				urls := make([]string, 0)
				for _, url := range community.(map[string]any) {
					if strings.TrimSpace(url.(string)) != `` {
						urls = append(urls, url.(string))
					}
				}
				dao.Socials[`social`] = urls
			}
			dao.ProjectImage = projectResearch.Image
			description, foundDescription := projectResearch.Detail[`description`]
			if foundDescription {
				dao.ProjectDescription = description.(string)
			} else {
				dao.ProjectDescription = ``
			}
			dao.ChainId = projectResearch.ChainId
			//Don't have chain id
			if strings.TrimSpace(projectResearch.ChainId) == `` {
				dao.ChainId = `0`
				dao.ChainName = `NULL`
			} else
			//Have chain id
			{
				chainList := &ChainList{}
				chainList.ChainId = projectResearch.ChainId
				chainList.SelectChainNameByChainId()
				//
				if strings.TrimSpace(chainList.ChainName) == `` {
					dao.ChainName = `NULL`
				} else {
					dao.ChainName = chainList.ChainName
				}
			}

			sourceCode, foundSourceCode := projectResearch.Detail[`sourceCode`]
			if foundSourceCode {
				dao.ExtraData[`sourceCode`] = sourceCode
			}

			website, foundWebsite := projectResearch.Detail[`website`]
			if foundWebsite {
				dao.ExtraData[`website`] = website
			}

			dao.ExtraData = make(map[string]any) //
			dao.Src = `researcher`
			dao.Tvl = 0
			dao.TotalUsed = 0
			projectRepo.Projects = append(projectRepo.Projects, dao)

			deletedIds = append(deletedIds, projectResearch.Id)
		}
	}

	fmt.Println(deletedIds)
}
