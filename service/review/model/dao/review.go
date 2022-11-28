package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/utils"
	dto "review-service/service/review/model/dto/revain"
)

type Review struct {
	Id          uint64
	AccountId   uint64
	Content     string
	Star        float64
	ProductId   uint64
	CreatedDate string
	UpdatedDate string
}

func (dao *Review) ConvertFrom(dto dto.ProductReview) {
	dao.ProductId = dto.ProductId
	dao.Content = dto.Content
	dao.Star = dto.Star
	dao.CreatedDate = dto.ReviewDate
}

func (dao *Review) InsertDB() error {
	query := `
	INSERT INTO review
		(accountid, "content", star, 
		productid, createddate, updateddate)
	VALUES
		( $1, $2, $3,
			$4, $5, $6)
		
		RETURNING id;
	`

	//Default value
	// dao.CreatedDate = utils.Timestamp() //override by crawled data
	dao.UpdatedDate = utils.Timestamp()

	var reviewId uint64
	err := db.PSQL.QueryRow(query,
		dao.AccountId, dao.Content, dao.Star,
		dao.ProductId, dao.CreatedDate, dao.UpdatedDate).Scan(&reviewId)
	dao.Id = reviewId

	return err
}
