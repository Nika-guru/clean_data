package dao

import (
	"encoding/json"
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"time"
)

type InvestorRepo struct {
	Investors []Investor
}

type Investor struct {
	InvestorCode  string
	InvestorName  string
	InvestorImage string
	CategoryName  string
	YearFounded   int
	Location      string
	Socials       map[string]any //map[string]string
	Description   string
	Src           string
	CreatedDate   string
	UpdatedDate   string
}

func (repo *InvestorRepo) InsertDB() {
	for _, dao := range repo.Investors {
		isExist, err := dao.IsExist()
		if err != nil {
			log.Println(log.LogLevelError, "service/review/model/dao/investor.go/(repo InvestorRepo) InsertDB()/dao.IsExist()", err.Error())
		}
		if !isExist {
			err := dao.InsertDB()
			if err != nil {
				time.Sleep(10 * time.Second)
				log.Println(log.LogLevelError, "service/review/model/dao/investor.go/(repo InvestorRepo) InsertDB()/dao.InsertDB()", err.Error())
			}
		}

	}
}
func (dao *Investor) IsExist() (isExist bool, err error) {
	query := `
		SELECT *
		FROM investor
		WHERE investorCode = $1
	`
	rows, err := db.PSQL.Query(query,
		dao.InvestorCode)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}
func (dao *Investor) InsertDB() error {
	query :=
		`
		INSERT INTO investor
			(investorcode, investorname, investorimage,
				categoryname, yearFounded, location,
				socials, description, src,
				createddate, updateddate)
		VALUES
			($1, $2, $3,
				$4, $5, $6,
				$7, $8, $9,
				$10, $11);
		
	`

	//default
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	websitesJSON, err := json.Marshal(dao.Socials)
	if err != nil {
		return err
	}

	_, err = db.PSQL.Exec(query,
		dao.InvestorCode, dao.InvestorName, dao.InvestorImage,
		dao.CategoryName, dao.YearFounded, dao.Location,
		websitesJSON, dao.Description, dao.Src,
		dao.CreatedDate, dao.UpdatedDate)
	return err
}
