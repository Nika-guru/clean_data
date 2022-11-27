package dao

import (
	"fmt"
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
	uniqueCategory := make(map[string]bool, 0)
	uniqueSubCategory := make(map[string]bool, 0)
	for _, dtoProduct := range dtoRepo.Products {
		for _, dtoProductCategory := range dtoProduct.ProductCategories {
			dtoProductCategory = strings.Trim(dtoProductCategory, ` `)
			trueVal, isCategory := constant.MAP_CATEGORY_PRODUCT_REVAIN[dtoProductCategory]

			if isCategory && trueVal {

				trueVal, foundInMap := uniqueCategory[dtoProductCategory]
				if foundInMap && trueVal {
				} else
				//search db
				{
					category := Category{}
					category.CategoryName = dtoProductCategory
					foundInDb, err := category.SelectByName()

					if err != nil {
						return err
					} else {
						if foundInDb {
						} else
						//Not found in db
						{
							fmt.Println(`====not found`, category.CategoryName)
							err := category.InsertDB()
							if err != nil {
								continue
							}
						}
						daoRepo.ProductCategories = append(daoRepo.ProductCategories, ProductCategory{
							ProductId:     *dtoProduct.ProductId,
							CategoryId:    category.CategoryId,
							SubCategoryId: nil,
							CreatedDate:   utils.Timestamp(),
							UpdatedDate:   utils.Timestamp(),
						})
						uniqueCategory[dtoProductCategory] = true
					}
				}
			} else
			//sub category
			{
				trueVal, foundInMap := uniqueSubCategory[dtoProductCategory]
				if foundInMap && trueVal {
				} else
				//search db
				{
					subCategory := SubCategory{}
					subCategory.SubCategoryName = dtoProductCategory
					foundInDb, err := subCategory.SelectByName()
					if err != nil {
						return err
					} else {
						if foundInDb {
						} else
						//Not found in db
						{
							err := subCategory.InsertDB()
							if err != nil {
								continue
							}
						}
						subCategoryIdVal := subCategory.SubCategoryId
						daoRepo.ProductCategories = append(daoRepo.ProductCategories, ProductCategory{
							ProductId:     *dtoProduct.ProductId,
							CategoryId:    subCategory.CategoryId,
							SubCategoryId: &subCategoryIdVal,
							CreatedDate:   utils.Timestamp(),
							UpdatedDate:   utils.Timestamp(),
						})
						uniqueSubCategory[dtoProductCategory] = true
					}
				}
			}
		}
	}
	return nil
}

func (daoRepo *ProductCategoryRepo) InsertDB() {
	for _, product := range daoRepo.ProductCategories {
		err := product.InsertDB()
		if err != nil {
			log.Println(log.LogLevelDebug, `service/reivew/model/dao/product_category.go/func (daoRepo *ProductCategoryRepo) ConverFrom(dtoRepo *dto.ProductInfoRepo) error/ err := product.InsertDB()`, err.Error())
			continue
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
