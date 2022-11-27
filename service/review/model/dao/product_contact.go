package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto "review-service/service/review/model/dto/revain"
)

type ProductContactRepo struct {
	ProductContacts []ProductContact
}

type ProductContact struct {
	ProductId   uint64
	Url         string
	Type        string
	CreatedDate string
	UpdatedDate string
}

func (repo *ProductContactRepo) InsertDB() {
	for _, product := range repo.ProductContacts {
		if !product.IsExist() {
			err := product.InsertDB()
			if err != nil {
				log.Println(log.LogLevelDebug, `service/reivew/model/dao/product_contact.go/func (repo *ProductContactRepo) InsertDB()/ product.InsertDB()`, err.Error())
				continue
			}
		}
	}
}

func (dao *ProductContact) IsExist() bool {
	query := `
		SELECT *
		FROM product_contact
		WHERE productid = $1 AND url = $2 AND "type" = $3;
	`
	rows, err := db.PSQL.Query(query,
		dao.ProductId, dao.Url, dao.Type)
	if err != nil {
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func (dao *ProductContact) InsertDB() error {
	query := `
	INSERT INTO product_contact
		(productid, url, "type",
		createddate, updateddate)
	VALUES
		($1, $2, $3,
		$4, $5);

	`

	//Set default value
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	_, err := db.PSQL.Exec(query,
		dao.ProductId, dao.Url, dao.Type,
		dao.CreatedDate, dao.UpdatedDate)
	return err
}

func (daoRepo *ProductContactRepo) ConvertFrom(dto dto.ProductDetail) {

	for url, typee := range dto.Url {
		daoRepo.ProductContacts = append(daoRepo.ProductContacts, ProductContact{
			ProductId: dto.ProductId,
			Url:       url,
			Type:      typee,
		})
	}

}
