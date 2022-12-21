package dao

import (
	"crawler/pkg/db"
	"crawler/pkg/log"
	"crawler/pkg/utils"
	"encoding/json"
	"fmt"
	"strings"
)

type ProductRepo struct {
	Products []*Product
}

type Product struct {
	Id           int //bigserial
	Type         string
	Address      string
	ChainId      string
	ChainName    string
	Name         string
	Image        string
	Description  string
	Category     string
	Subcategory  string
	Detail       map[string]any
	CreatedDate  string
	UpdatedDate  string
	TotalReviews uint64
	TotalIsScam  uint64
	TotalNotScam uint64
	IsScam       bool
	FromBy       string
	Contract     map[string]any
	Symbol       string
	IsShow       bool
	TotalWarning uint64
	Holders      uint64
	Transfers    uint64
	SourceCode   string
	ContractABI  string
}

func (repo *ProductRepo) InsertDB(sources map[string]bool) {
	for _, product := range repo.Products {
		isExist, err := product.IsExist1(sources)
		if err != nil {
			log.Println(log.LogLevelError, `service/merge/model/dao/product.go/ (repo *ProductRepo) InsertDB()/ product.IsExist(), at product name `+product.Name+` from src: `+product.FromBy, err.Error())
			continue
		}
		if !isExist {
			err = product.InsertDB()
			if err != nil {
				log.Println(log.LogLevelError, `service/merge/model/dao/product.go/ (repo *ProductRepo) InsertDB()/ product.InsertDB(), at product name  `+product.Name+` from src: `+product.FromBy, err.Error())
				for source := range sources {
					productCode := product.Detail[`productCode`+strings.Title(strings.TrimSpace(source))].(string)
					product.InsertFail(productCode, source, err.Error())
				}
				continue
			}
		}
	}
}

func (dao *Product) IsExist1(sources map[string]bool) (isExist bool, err error) {
	query :=
		`
		SELECT 
			*
		FROM
			product
		WHERE 

		`
	commonSource := ``
	idx := 0
	for source := range sources {
		commonSource = fmt.Sprintf("productCode%s", strings.Title(strings.ToLower(source)))
		query += fmt.Sprintf("detail->>'%s' = $1", commonSource)

		//don't append reduant 'or' at query
		if len(sources)-1 != idx {
			query += ` or `
		}
		idx++
	}

	rows, err := db.PSQL.Query(query, dao.Detail[commonSource])
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (dao *Product) IsExist(sources map[string]bool, code string) (isExist bool, err error) {
	query :=
		`
		SELECT 
			*
		FROM
			product
		WHERE 

		`
	commonSource := ``
	idx := 0
	for source := range sources {
		commonSource = fmt.Sprintf("productCode%s", strings.Title(strings.ToLower(source)))
		query += fmt.Sprintf("detail->>'%s' = $1", commonSource)

		//don't append reduant 'or' at query
		if len(sources)-1 != idx {
			query += ` or `
		}
		idx++
	}

	rows, err := db.PSQL.Query(query, code)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (dao *Product) InsertDB() error {
	query := `
	INSERT INTO product
		("type", address, chainid,
		chainname, "name", image,
		"desc", category, subcategory,
		detail, createddate, updateddate,
		totalreviews, totalisscam, totalnotscam,
		isscam, fromby, contract,
		symbol, isshow)
	VALUES
		($1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			$13, $14, $15,
			$16, $17, $18,
			$19, $20)
	RETURNING id;
	`
	//Default value
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()
	dao.TotalReviews = 0
	dao.TotalIsScam = 0
	dao.TotalNotScam = 0
	dao.IsScam = false
	dao.IsShow = true //enough data

	detailJSON, err := json.Marshal(dao.Detail)
	if err != nil {
		return err
	}

	contractJSON, err := json.Marshal(dao.Contract)
	if err != nil {
		return err
	}

	err = db.PSQL.QueryRow(query,
		dao.Type, dao.Address, dao.ChainId,
		dao.ChainName, dao.Name, dao.Image,
		dao.Description, dao.Category, dao.Subcategory,
		detailJSON, dao.CreatedDate, dao.UpdatedDate,
		dao.TotalReviews, dao.TotalIsScam, dao.TotalNotScam,
		dao.IsScam, dao.FromBy, contractJSON,
		dao.Symbol, dao.IsShow).Scan(&dao.Id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductRepo) FormatCategoryDappRadar() {
	for _, product := range repo.Products {
		switch product.Category {
		case `exchanges`:
			fallthrough
		case `defi`:
			product.Category = `Crypto Exchanges`
		case `games`:
			product.Category = `Blockchain Games`
		case `marketplaces`:
			product.Category = `NFT Marketplaces`
		default:
			product.Category = strings.Title(strings.ToLower(product.Category))
		}
	}
}

func (repo *ProductRepo) FormatProductNameRevain() {
	for _, product := range repo.Products {
		product.Name = strings.Title(strings.ToLower(product.Name))
		product.Name = strings.ReplaceAll(product.Name, `-`, ` `)
	}
}

//#################### Info ###################
func (dao *Product) SelectTotal(source string) (uint64, error) {
	query :=
		`
		SELECT COUNT(*) FROM product WHERE detail->>'productCode` + strings.Title(strings.TrimSpace(source)) + `' IS NOT NULL;
	`
	rows, err := db.PSQL.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total uint64 = 0
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (dao *Product) SelectTotalIn24h(source string) (uint64, error) {
	query :=
		`
		SELECT COUNT(*) FROM product 
		WHERE detail->>'productCode` + strings.Title(strings.TrimSpace(source)) + `' IS NOT NULL
			AND createddate BETWEEN (NOW() - INTERVAL '24 HOURS') AND NOW();
	`
	rows, err := db.PSQL.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total uint64 = 0
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (dao *Product) SelectLatestCreatedDate(source string) (string, error) {
	query := `
		SELECT createddate 
		FROM product
		WHERE detail->>'productCode` + strings.Title(strings.TrimSpace(source)) + `' IS NOT NULL 
		ORDER BY createddate DESC
		LIMIT 1
	`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return ``, err
	}
	defer rows.Close()

	var latestCreatedDate string = ``
	if rows.Next() {
		err = rows.Scan(&latestCreatedDate)
		if err != nil {
			return ``, err
		}
	}

	return latestCreatedDate, nil
}

func (dao *Product) SelectCrawledFailedTotal(source string) (uint64, error) {
	query := `
		SELECT COUNT(*) 
		FROM crawl_fail_3rd 
		WHERE src = '` + strings.TrimSpace(source) + `'`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total uint64 = 0
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (dao *Product) InsertFail(productCode string, source string, errMsg string) error {
	query := `
		INSERT INTO crawl_fail_3rd
			(productcode, src, err)
		VALUES
			($1, $2, $3);
	`
	_, err := db.PSQL.Exec(query,
		productCode, source, errMsg)
	return err
}
