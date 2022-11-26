package dao

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto "review-service/service/review/model/dto/revain"
)

type ProductRepo struct {
	Products []Product
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

func (repo *ProductRepo) InsertDB() {
	for _, product := range repo.Products {
		err := product.InsertDB()
		if err != nil {
			//TODO: insert fail success + sleep
			log.Println(log.LogLevelDebug, `service/reivew/model/dao/product_info.go/func (repo *ProductInfoRepo) InsertDB()/ productInfo.InsertDB()`, err.Error())
			continue
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

	err = db.PSQL.QueryRow(query,
		product.ProductName, product.ProductImage, product.ProductDescription,
		productDetailJSONB, product.CrawlSource, product.CreatedDate,
		product.UpdatedDate).Scan(&product.ProductId)
	return err
}

func (daoRepo *ProductRepo) ConverFrom(dtoRepo *dto.ProductInfoRepo) {

	for _, productDto := range dtoRepo.Products {
		productDao := Product{
			ProductName:   productDto.ProductName,
			ProductImage:  productDto.ProductImage,
			ProductDetail: productDto.ProductDetail,
			CrawlSource:   productDto.CrawlSource,
		}
		daoRepo.Products = append(daoRepo.Products, productDao)
	}
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
		product := Product{}
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
