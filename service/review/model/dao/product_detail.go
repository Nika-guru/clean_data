package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
)

type ProductDetailRepo struct {
	ProductDetailList []ProductDetail
}

type ProductDetail struct {
	ProductId       int64
	Description     string
	OfficialWebsite string
	SocialMedia     string
	CreatedDate     string
	UpdatedDate     string
}

func (repo *ProductDetailRepo) InsertDB(productInfoMap map[int64]Product) {
	for _, productDetail := range repo.ProductDetailList {
		// Get a Tx for making transaction requests.
		tx, err := db.PSQL.Begin()
		if err != nil {
			//TODO: insert fail_model
			continue
		}

		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
			err = tx.Commit()
		}()

		err = productDetail.InsertDB()
		if err != nil {
			//TODO: insert fail success + sleep
			log.Println(log.LogLevelDebug, `service/reivew/model/dao/product_detail.go/ (repo *ProductDetailRepo) InsertDB(productInfoMap map[int64]ProductInfoRepo)/ productDetail.InsertDB()`, err.Error())
			continue
		}
	}
}

func (dao *ProductDetail) InsertDB() error {
	query := `
	INSERT INTO product_detail
		(productid, description, officialwebsite,
			socialmedia, createddate, updateddate)
	VALUES($1, $2, $3,
			$4, $5, $6);
	`

	//Set default value
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	_, err := db.PSQL.Query(query,
		dao.ProductId, dao.Description, dao.OfficialWebsite,
		dao.SocialMedia, dao.CreatedDate, dao.UpdatedDate)

	if err != nil {
		return err
	}
	return nil
}
