package dao

import (
	"base/pkg/db"
	"encoding/json"
	"strings"
)

type ProjectRepo struct {
	Projects []*Project
}

type Project struct {
	ProjectCode        string //id
	ProjectName        string
	ProjectCategory    string
	ProjectSubcategory string
	Socials            map[string]any
	ProjectImage       string
	ProjectDescription string
	ChainId            string
	ChainName          string
	ExtraData          map[string]any
	CreatedDate        string
	UpdatedDate        string
	Src                string
	Tvl                float64
	TotalUsed          float64
}

type ContractData struct {
	ChainId   string `json:"chainId"`
	ChainName string `json:"chainName"`
	Address   string `json:"address"`
}

func (repo *ProjectRepo) SelectBySrc(src string) error {
	query :=
		`
	SELECT 
		id, "name", category,
		subcategory, social, image,
		description, chainid, chainname,
		extradata, createddate, updateddate,
		src, tvl, totalused
	FROM 
		project
	WHERE
		src = $1
	`

	rows, err := db.PSQL.Query(query, src)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		dao := &Project{}
		socialsJSON := []byte{}
		extraDataJSON := []byte{}

		rows.Scan(&dao.ProjectCode, &dao.ProjectName, &dao.ProjectCategory,
			&dao.ProjectSubcategory, &socialsJSON, &dao.ProjectImage,
			&dao.ProjectDescription, &dao.ChainId, &dao.ChainName,
			&extraDataJSON, &dao.CreatedDate, &dao.UpdatedDate,
			&dao.Src, &dao.Tvl, &dao.TotalUsed)

		if socialsJSON != nil {
			err := json.Unmarshal(socialsJSON, &dao.Socials)
			if err != nil {
				return err
			}
		}
		if extraDataJSON != nil {
			err = json.Unmarshal(extraDataJSON, &dao.ExtraData)
			if err != nil {
				return err
			}
		}

		repo.Projects = append(repo.Projects, dao)
	}

	return nil
}

func (dao *Project) UpdateSocialsAndExtraDataByCodeAndSrc() error {
	query :=
		`
		UPDATE 
			project
		SET 
			social=$1,
			extradata =$4
		WHERE 
			src = $3
		AND
			id = $2
	`
	socialsJSON, err := json.Marshal(dao.Socials)
	if err != nil {
		return err
	}
	extraDataJSON, err := json.Marshal(dao.ExtraData)
	if err != nil {
		return err
	}

	_, err = db.PSQL.Exec(query, socialsJSON, dao.ProjectCode, dao.Src, extraDataJSON)
	if err != nil {
		return err
	}

	return nil
}

func (dao *Project) IsUpdatedSocialsByCodeAndSrc() (isUpdated bool, err error) {
	query :=
		`
		SELECT *
		FROM tmp_updated_socials_dapp
		WHERE project_code = $1 AND src = $2
	`

	rows, err := db.PSQL.Query(query, dao.ProjectCode, dao.Src)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (dao *Project) InsertUpdatedSocial() error {
	query := `
			INSERT INTO tmp_updated_socials_dapp
		(project_code, src)
		VALUES($1, $2);
	`
	_, err := db.PSQL.Exec(query, dao.ProjectCode, dao.Src)
	return err
}

func (repo *ProjectRepo) SelectMergeCommonDataFromDappAndRevain() error {
	query :=
		`
		SELECT 
			dappradar.id, dappradar."name", revan.category,
			dappradar.subcategory, dappradar.social, dappradar.image,
			dappradar.description, chainList.chainId, dappradar.chainname,
			dappradar.extradata, dappradar.createddate, dappradar.updateddate,
			'3rd' as src, dappradar.tvl, dappradar.totalused
		FROM 
				(SELECT * FROM project WHERE src = 'dappradar') AS dappradar
			INNER JOIN
				(SELECT * FROM project WHERE src = 'revan') AS revan
				ON revan.name = dappradar.id
			LEFT JOIN
				(SELECT * FROM chain_list) AS chainList
				ON chainlist.chainname = dappradar.chainname
		WHERE 
			(	(dappradar.category = 'defi' AND revan.category = 'Crypto Exchanges')
			OR
				(dappradar.category = 'exchanges' AND revan.category = 'Crypto Exchanges')
			OR
				(dappradar.category = 'games' AND revan.category = 'Blockchain Games')
			OR
				(dappradar.category = 'marketplaces' AND revan.category = 'NFT Marketplaces')
			)
		ORDER BY dappradar.id, chainid
	`
	return repo.Select(query)
}

func (repo *ProjectRepo) Select(query string) error {
	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		dao := &Project{}
		socialsJSON := []byte{}
		extraDataJSON := []byte{}
		var subcategory any //can be null when absent data
		var chainName any   //can be null when absent data
		var chainId any     //can be null when left join

		rows.Scan(
			&dao.ProjectCode, &dao.ProjectName, &dao.ProjectCategory,
			&subcategory, &socialsJSON, &dao.ProjectImage,
			&dao.ProjectDescription, &chainId, &chainName,
			&extraDataJSON, &dao.CreatedDate, &dao.UpdatedDate,
			&dao.Src, &dao.Tvl, &dao.TotalUsed)

		if len(socialsJSON) != 0 {
			err := json.Unmarshal(socialsJSON, &dao.Socials)
			if err != nil {
				return err
			}
		}

		if len(extraDataJSON) != 0 {
			err := json.Unmarshal(extraDataJSON, &dao.ExtraData)
			if err != nil {
				return err
			}
		}

		if subcategory == nil {
			dao.ProjectSubcategory = ``
		} else {
			dao.ProjectSubcategory = strings.TrimSpace(subcategory.(string))
		}

		if chainId == nil {
			dao.ChainId = `0`
		} else {
			dao.ChainId = chainId.(string)
		}

		if chainName == nil {
			dao.ChainName = `NULL`
		} else {
			dao.ChainName = strings.TrimSpace(chainName.(string))
		}

		repo.Projects = append(repo.Projects, dao)
	}
	return nil
}
func (repo *ProjectRepo) ConvertToProductRepo(productRepo *ProductRepo, sources map[string]bool) {
	oldProjectCode := ``

	contracts := make(map[string]any, 0)
	contracts[`contract`] = make([]ContractData, 0)
	product := &Product{}
	for _, project := range repo.Projects {
		if oldProjectCode != project.ProjectCode {
			contracts := make(map[string]any, 0)
			contracts[`contract`] = make([]ContractData, 0)

			product = &Product{}
			project.ConvertToProduct(product, sources)
			productRepo.Products = append(productRepo.Products, product)
		}

		address, foundAddress := project.ExtraData[`address`]
		if !foundAddress {
			address = `NULL`
		}

		if project.ChainName != `NULL` && address.(string) != `NULL` {
			contracts[`contract`] = append(contracts[`contract`].([]ContractData), ContractData{
				ChainId:   project.ChainId,
				ChainName: project.ChainName,
				Address:   address.(string),
			})
		}
		oldProjectCode = project.ProjectCode

		product.Contract = contracts //append update
	}
}

func (project *Project) ConvertToProduct(product *Product, sources map[string]bool) {
	product.Type = `project`
	//Non-EVM
	if strings.TrimSpace((project.ChainId)) == `` || strings.TrimSpace((project.ChainId)) == `0` || strings.TrimSpace((project.ChainId)) == `NULL` || strings.TrimSpace((project.ChainId)) == `undefined` {
		product.ChainId = `0`
	} else {
		product.ChainId = project.ChainId
	}
	if strings.TrimSpace(project.ChainName) == `` {
		product.ChainName = `NULL`
	} else {
		product.ChainName = project.ChainName
	}
	product.Name = project.ProjectName
	product.Image = project.ProjectImage
	product.Description = project.ProjectDescription
	product.Category = project.ProjectCategory
	product.Subcategory = project.ProjectSubcategory
	if product.Detail == nil {
		product.Detail = make(map[string]any, 0)
	}
	val, foundProject := project.ExtraData[`address`]
	if foundProject {
		product.Address = strings.ToLower(val.(string))
	} else {
		product.Address = `NULL`

	}

	_, foundProduct := product.Detail[`tvl`]
	if !foundProduct {
		product.Detail[`tvl`] = project.Tvl
	}
	_, foundProduct = product.Detail[`totalUsed`]
	if !foundProduct {
		product.Detail[`totalUsed`] = project.TotalUsed
	}

	_, foundProduct = product.Detail[`social`]
	_, foundProject = project.Socials[`social`]
	if !foundProduct && foundProject {
		product.Detail[`social`] = project.Socials[`social`]
	}

	_, foundProduct = product.Detail[`website`]
	_, foundProject = project.ExtraData[`website`]
	if !foundProduct && foundProject {
		product.Detail[`website`] = project.ExtraData[`website`]
	}

	_, foundProduct = product.Detail[`sourceCode`]
	_, foundProject = project.ExtraData[`sourceCode`]
	if !foundProduct && foundProject {
		product.Detail[`sourceCode`] = project.ExtraData[`sourceCode`]
	}

	for source := range sources {
		key := `productCode` + strings.Title(strings.ToLower(source))
		_, foundProduct = product.Detail[key]
		product.Detail[key] = project.ProjectCode
	}

	product.FromBy = project.Src
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func (repo *ProjectRepo) SelectDappradarWithoutRevain() error {
	query := `
	--only revain exist, dappradar don't have any data
	SELECT 
		project.id, project."name", project.category,
		project.subcategory, project.social, project.image,
		project.description, chain_list.chainId, project.chainname,
		project.extradata, project.createddate, project.updateddate,
		'3rd' as src, project.tvl, project.totalused
	FROM 
		(SELECT *	
		FROM project
		WHERE 
			id NOT IN (SELECT * FROM tmp_project_id_dappradar_intersect_revain)
			AND
			src = 'dappradar'
		) AS project
	LEFT JOIN chain_list
		ON project.chainname  = chain_list.chainname 
	ORDER BY id, chainId;
	`

	return repo.Select(query)
}

func (repo *ProjectRepo) SelectRevainWithoutDappradar() error {
	query := `
	--only revain exist, dappradar don't have any data
	SELECT 
		project."name" as id, project."name", project.category,
		project.subcategory, project.social, project.image,
		project.description, chain_list.chainId, project.chainname,
		project.extradata, project.createddate, project.updateddate,
		'3rd' as src, project.tvl, project.totalused
	FROM 
		(SELECT *	
		FROM project
		WHERE 
			"name" NOT IN (SELECT * FROM tmp_project_id_dappradar_intersect_revain)
			AND
			src = 'revain'
		) AS project
	LEFT JOIN chain_list
		ON project.chainname  = chain_list.chainname 
	ORDER BY "name", chainId
	`

	return repo.Select(query)
}

func (repo *ProjectRepo) SelectProjectSearch() error {
	query := `
	`

	return repo.Select(query)
}
