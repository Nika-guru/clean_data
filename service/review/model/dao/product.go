package dao

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto "review-service/service/review/model/dto/revain"
	dto_revain "review-service/service/review/model/dto/revain"
)

type ProductRepo struct {
	Products []*Product
}

type Product struct {
	ProductId          uint64         `json:"product_id"`
	ProductName        string         `json:"title"`
	ProductImage       string         `json:"image"`
	ProductDescription string         `json:"description"`
	ProductDetail      map[string]any `json:"detail"`
	CrawlSource        string         `json:"crawl_source"`
	CreatedDate        string         `json:"created_date"`
	UpdatedDate        string         `json:"updated_date"`
}

func (repo *ProductRepo) InsertDB(productInfoDebug *dto_revain.ProductInfoDebug) {
	for _, product := range repo.Products {
		//Not exist, got existed id(incremental)
		isExist, err := product.SelectByProductName()
		if err != nil {
			log.Println(log.LogLevelDebug, `service/reivew/model/dao/product.go/func (repo *ProductInfoRepo) InsertDB()/ product.SelectByProductName()`, err.Error())
			continue
		}
		if !isExist {
			err := product.InsertDB()
			if err != nil {
				log.Println(log.LogLevelDebug, `service/reivew/model/dao/product_info.go/func (repo *ProductInfoRepo) InsertDB()/ productInfo.InsertDB()`, err.Error())
				continue
			}
			//######Start Debuging #########
			debug := dto_revain.Debug{}
			productInfoDebug.IsSuccess = true
			debug.AddProductInfo(*productInfoDebug)
			//######End Debuging #########
		}
	}
}

func (product *Product) InsertDB() error {
	query := `
	INSERT INTO product
		(productname, productimage, productdescription,
			productdetail, crawlsource, createddate,
			updateddate)
	VALUES
		($1 , $2 , $3 ,
			$4 , $5 , $6 ,
			$7)
		RETURNING productid;
	`

	//Set default value
	product.CreatedDate = utils.Timestamp()
	product.UpdatedDate = utils.Timestamp()

	productDetailJSONB, err := json.Marshal(product.ProductDetail)
	if err != nil {
		return err
	}

	var id uint64
	err = db.PSQL.QueryRow(query,
		product.ProductName, product.ProductImage, product.ProductDescription,
		productDetailJSONB, product.CrawlSource, product.CreatedDate,
		product.UpdatedDate).Scan(&id)
	product.ProductId = id
	return err
}

func (daoRepo *ProductRepo) ConverFrom(dtoRepo *dto.ProductInfoRepo) {

	for _, productDto := range dtoRepo.Products {
		productDao := &Product{
			ProductName:   productDto.ProductName,
			ProductImage:  productDto.ProductImage,
			ProductDetail: productDto.ProductDetail,
			CrawlSource:   productDto.CrawlSource,
		}
		productDto.ProductId = &productDao.ProductId
		daoRepo.Products = append(daoRepo.Products, productDao)
	}
}

func (dao *Product) ConvertFrom(dto dto.ProductDetail) {
	dao.ProductId = dto.ProductId
	dao.ProductDescription = dto.Description
}

func (dao *Product) UpdateDescription() error {
	query := `
	UPDATE product
		SET productdescription=$1
		WHERE productid = $2
	`
	_, err := db.PSQL.Exec(query, dao.ProductDescription, dao.ProductId)
	if err != nil {
		return err
	}

	return nil
}

func (dao *Product) SelectByProductName() (isExist bool, err error) {
	query := `
		SELECT DISTINCT ON(productname) 
			productid, productname, productimage,
			productdescription, productdetail, crawlsource,
			createddate, updateddate
		FROM product
		WHERE productname = $1
		ORDER BY productname, productid; --get smallest id if duplicate
	`
	rows, err := db.PSQL.Query(query, dao.ProductName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		productDetailJSONB := []byte{}
		err := rows.Scan(
			&dao.ProductId, &dao.ProductName, &dao.ProductImage,
			&dao.ProductDescription, &productDetailJSONB, &dao.CrawlSource,
			&dao.CreatedDate, &dao.UpdatedDate)
		if err != nil {
			return true, err
		}
		json.Unmarshal(productDetailJSONB, &dao.ProductDetail)
		if err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func (repo *ProductRepo) SelectByTitleWithShortDescriptionAndType(keyword string, productTypeId int64) error {
	//DISTINCT for avoiding duplicate same product id same name, same group it belongs to
	query := `
	SELECT
	DISTINCT ON (title, producttypenames) 
		productid, title, producttypename,
		image, shortdescription, averagestar,
		totalreviews, createddate, updateddate, 
		producttypenames
	FROM product_info
	WHERE
		(
		LOWER(title) LIKE LOWER('%` + keyword + `%')
		OR LOWER(shortdescription) LIKE LOWER('%` + keyword + `%')
		)
		AND 
	`
	//for all type
	if productTypeId == -1 {
		query += ` 1 = 1 `
	} else
	//specified type
	{
		//integer, no need valudate
		query += fmt.Sprintf(` producttypename = (select "name" from product_type where id = %d ) `, productTypeId)
	}
	//order to review latest review post (smallest productid)
	query += ` ORDER BY title, producttypenames, productid `

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		product := &Product{}
		productDetailJSONB := []byte{}
		rows.Scan(&product.ProductName, &product.ProductImage, &product.ProductDescription,
			&productDetailJSONB, &product.CrawlSource, &product.CreatedDate,
			&product.UpdatedDate)
		err := json.Unmarshal(productDetailJSONB, &product.ProductDetail)
		if err != nil {
			return err
		}
		repo.Products = append(repo.Products, product)
	}

	return nil
}

func (product *Product) IsExist() bool {
	query := `
		SELECT *
		FROM product
		WHERE productname = $1 AND productimage = $2 AND crawlsource = $3;
	`
	rows, err := db.PSQL.Query(query,
		product.ProductName, product.ProductImage, product.CrawlSource)
	if err != nil {
		return false
	}
	defer rows.Close()

	return rows.Next()
}
