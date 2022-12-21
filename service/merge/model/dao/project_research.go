package dao

import (
	"crawler/pkg/db"
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

func (repo *ProjectResearchRepo) ConvertToProductRepo(productRepo *ProductRepo) {
	deletedIds := make([]string, 0)
	for _, projectResearch := range repo.ProjectResearchList {
		switch strings.TrimSpace(projectResearch.Type) {
		//Coin/ token/ ico
		case `coin`:
			fallthrough
		case `token`:
			fallthrough
		case `ico`:
			dao := &Product{}

			dao.Type = projectResearch.Type
			if projectResearch.Address == `` {
				dao.Address = `NULL`
			} else {
				dao.Address = projectResearch.Address
			}
			if projectResearch.ChainId == `` {
				dao.ChainId = `NULL`
				dao.ChainName = `NULL`
			} else {
				dao.ChainId = projectResearch.ChainId

				chainList := ChainList{}
				chainList.ChainId = projectResearch.ChainId
				chainList.SelectChainNameByChainId()
				if chainList.ChainName == `` {
					dao.ChainName = `NULL`
				} else {
					dao.ChainName = chainList.ChainName
				}
			}
			dao.Symbol = strings.ToUpper(projectResearch.Symbol)
			dao.Name = projectResearch.Name
			dao.Category = projectResearch.Category
			dao.Subcategory = projectResearch.Subcategory
			dao.Detail = make(map[string]any)

			desc, foundDesc := projectResearch.Detail[`description`]
			if foundDesc {
				dao.Description = desc.(string)
			}
			dao.Detail = projectResearch.Detail
			dao.FromBy = `research`
			dao.Detail[`totalSupply`] = 0
			dao.Detail[`maxSupply`] = projectResearch.TotalSupply
			dao.Detail[`marketcap`] = 0
			dao.Detail[`volumeTrading`] = 0
			dao.Detail[`holder`] = 0
			dao.Detail[`decimals`] = projectResearch.Decimals

			dao.Image = projectResearch.Image

			productRepo.Products = append(productRepo.Products, dao)

			deletedIds = append(deletedIds, projectResearch.Id)
			//dApp
		case `project`:
			// dao := &Project{}
			// dao.ProjectCode = ``
			// dao.ProjectName = projectResearch.Name
			// dao.ProjectCategory = projectResearch.Category
			// dao.ProjectSubcategory = projectResearch.Subcategory

			// //avoid nil pointer below
			// if projectResearch.Detail == nil {
			// 	projectResearch.Detail = make(map[string]any, 0)
			// }

			// community, foundCommunity := projectResearch.Detail[`community`]
			// if foundCommunity && len(community.(map[string]any)) > 0 {
			// 	dao.Socials = make(map[string]any, 0)
			// 	dao.Socials[`social`] = make([]string, 0)
			// 	urls := make([]string, 0)
			// 	for _, url := range community.(map[string]any) {
			// 		if strings.TrimSpace(url.(string)) != `` {
			// 			urls = append(urls, url.(string))
			// 		}
			// 	}
			// 	dao.Socials[`social`] = urls
			// }
			// dao.ProjectImage = projectResearch.Image
			// description, foundDescription := projectResearch.Detail[`description`]
			// if foundDescription {
			// 	dao.ProjectDescription = description.(string)
			// } else {
			// 	dao.ProjectDescription = ``
			// }
			// dao.ChainId = projectResearch.ChainId
			// //Don't have chain id
			// if strings.TrimSpace(projectResearch.ChainId) == `` {
			// 	dao.ChainId = `0`
			// 	dao.ChainName = `NULL`
			// } else
			// //Have chain id
			// {
			// 	chainList := &ChainList{}
			// 	chainList.ChainId = projectResearch.ChainId
			// 	chainList.SelectChainNameByChainId()
			// 	//
			// 	if strings.TrimSpace(chainList.ChainName) == `` {
			// 		dao.ChainName = `NULL`
			// 	} else {
			// 		dao.ChainName = chainList.ChainName
			// 	}
			// }

			// sourceCode, foundSourceCode := projectResearch.Detail[`sourceCode`]
			// if foundSourceCode {
			// 	dao.ExtraData[`sourceCode`] = sourceCode
			// }

			// website, foundWebsite := projectResearch.Detail[`website`]
			// if foundWebsite {
			// 	dao.ExtraData[`website`] = website
			// }

			// dao.ExtraData = make(map[string]any) //
			// dao.Src = `research`
			// dao.Tvl = 0
			// dao.TotalUsed = 0
			// productRepo.Products = append(productRepo.Products, dao)

			// deletedIds = append(deletedIds, projectResearch.Id)
		}
	}

	fmt.Println(deletedIds)
}
