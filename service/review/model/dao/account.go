package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/utils"
	dto "review-service/service/review/model/dto/revain"
)

type Account struct {
	Id          uint64
	UserId      uint64 //fb, twitter
	AccountType string

	Role        uint8
	Password    string
	Image       string
	Email       string
	Username    string
	CreatedDate string
	UpdatedDate string
}

func (dao *Account) ConvertFrom(dto dto.ProductReview) {
	dao.Username = dto.Username
	dao.Image = dto.AccountImage
}

func (dao *Account) InsertDB() error {
	query := `
	INSERT INTO account
		(image, accounttype, username, 
			createddate, updateddate)
	VALUES
		( $1, $2, $3,
			$4, $5)
			
		RETURNING id;
	
	`

	//Default value
	dao.AccountType = `user`
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	var accountId uint64
	err := db.PSQL.QueryRow(query,
		dao.Image, dao.AccountType, dao.Username,
		dao.CreatedDate, dao.UpdatedDate).Scan(&accountId)
	dao.Id = accountId

	return err
}
