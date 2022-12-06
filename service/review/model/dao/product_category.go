package dao

import (
	"database/sql"
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	dto "review-service/service/review/model/dto/revain"
	"strings"
)

type ProductCategoryRepo struct {
	ProductCategories []ProductCategory
}

type ProductCategory struct {
	ProductId     uint64
	CategoryId    uint64
	SubCategoryId *uint64 //can Null
	CreatedDate   string
	UpdatedDate   string
}

func (daoRepo *ProductCategoryRepo) ConverFrom(dtoRepo *dto.ProductInfoRepo) error {
	uniqueCategory := make(map[string]Category, 0)
	uniqueSubCategory := make(map[string]SubCategory, 0)
	for _, dtoProduct := range dtoRepo.Products {
		for _, dtoProductCategory := range dtoProduct.ProductCategories {
			dtoProductCategory = strings.Trim(dtoProductCategory, ` `)
			trueVal, isCategory := constant.MAP_CATEGORY_PRODUCT_REVAIN[dtoProductCategory]

			if isCategory && trueVal {

				category, foundInMap := uniqueCategory[dtoProductCategory]
				if foundInMap {
				} else
				//search db
				{
					category = Category{}
					category.CategoryName = dtoProductCategory
					foundInDb, err := category.SelectByName()

					if err != nil {
						return err
					}

					if foundInDb {
					} else
					//Not found in db
					{
						err := category.InsertDB()
						if err != nil {
							//continue
							return err
						}
					}

					uniqueCategory[dtoProductCategory] = category
				}

				daoRepo.ProductCategories = append(daoRepo.ProductCategories, ProductCategory{
					ProductId:     *dtoProduct.ProductId,
					CategoryId:    category.CategoryId,
					SubCategoryId: nil,
					CreatedDate:   utils.Timestamp(),
					UpdatedDate:   utils.Timestamp(),
				})
			} else
			//sub category
			{
				subCategory, foundInMap := uniqueSubCategory[dtoProductCategory]
				if foundInMap {
				} else
				//search db
				{
					subCategory = SubCategory{}
					subCategory.SubCategoryName = dtoProductCategory
					foundInDb, err := subCategory.SelectByName()
					if err != nil {
						return err
					}
					if foundInDb {
					} else
					//Not found in db
					{
						category := Category{}
						category.CategoryName = constant.DEFAULT_CATEGORY_PRODUCT_REVAIN
						foundInDb, err := category.SelectByName()
						if err != nil {
							return err
						}
						if foundInDb {
							subCategory.CategoryId = category.CategoryId
						}
						err = subCategory.InsertDB()
						if err != nil {
							//continue
							return err
						}
					}
					uniqueSubCategory[dtoProductCategory] = subCategory
				}
				subCategoryIdVal := subCategory.SubCategoryId
				daoRepo.ProductCategories = append(daoRepo.ProductCategories, ProductCategory{
					ProductId:     *dtoProduct.ProductId,
					CategoryId:    subCategory.CategoryId,
					SubCategoryId: &subCategoryIdVal,
					CreatedDate:   utils.Timestamp(),
					UpdatedDate:   utils.Timestamp(),
				})
			}
		}
	}
	return nil
}

func (daoRepo *ProductCategoryRepo) InsertDB() {
	for _, productCategory := range daoRepo.ProductCategories {
		if !productCategory.IsExist() {
			err := productCategory.InsertDB()
			if err != nil {
				log.Println(log.LogLevelDebug, `service/reivew/model/dao/product_category.go/func (daoRepo *ProductCategoryRepo) InsertDB()/ err := productCategory.InsertDB()`, err.Error())
				continue
			}

			//ban bum bum
		}
	}
}

func (dao *ProductCategory) InsertDB() error {
	query := `
	INSERT INTO product_category
		(productid, categoryid, subcategoryid, 
			createddate, updateddate)
	VALUES
		($1, $2, $3,
			$4, $5);	
	`

	//Set default value
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	_, err := db.PSQL.Exec(query,
		dao.ProductId, dao.CategoryId, dao.SubCategoryId,
		dao.CreatedDate, dao.UpdatedDate)

	return err
}

func (productCategory *ProductCategory) IsExist() bool {
	query := `
		SELECT *
		FROM product_category
		WHERE productid = $1 AND categoryid = $2 AND subcategoryid
	`

	var rows *sql.Rows
	var err error
	if productCategory.SubCategoryId != nil {
		query += ` = $3;`
		rows, err = db.PSQL.Query(query,
			productCategory.ProductId, productCategory.CategoryId, *productCategory.SubCategoryId)
	} else {
		query += ` IS NULL`
		rows, err = db.PSQL.Query(query,
			productCategory.ProductId, productCategory.CategoryId)
	}

	if err != nil {
		return false
	}
	defer rows.Close()

	return rows.Next()
}
