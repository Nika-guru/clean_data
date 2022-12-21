package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"time"
)

type FuncRaisingRepo struct {
	FuncRaisingList []FuncRaising
}
type FuncRaising struct {
	ProjectCode string
	ProjectName string
	ProjectLogo string

	InvestorCode string
	InvestorName string
	InvestorLogo string

	FundStageCode string
	FundStageName string
	FundAmount    float64
	FundDate      string

	CreatedDate string
	UpdatedDate string

	Description     string
	AnnouncementUrl string

	Valulation  float64
	SrcFund     string
	SrcInvestor string
}

func (dao FuncRaising) IsExist() (isExist bool, err error) {
	query :=
		`SELECT * 
		FROM fund_raising 
		WHERE projectcode = $1 
			AND investorcode = $2 
			AND fundstagecode = $3`

	rows, err := db.PSQL.Query(query,
		dao.ProjectCode, dao.InvestorCode, dao.FundStageCode)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (repo FuncRaisingRepo) InsertDB() {
	for _, dao := range repo.FuncRaisingList {
		isExist, err := dao.IsExist()
		if err != nil {
			log.Println(log.LogLevelError, "service/review/model/dao/fund_raising.go/(repo FuncRaisingRepo) InsertDB()/dao.IsExist()", err.Error())
		}
		if !isExist {
			err := dao.InsertDB()
			if err != nil {
				time.Sleep(10 * time.Second)
				log.Println(log.LogLevelError, "service/review/model/dao/fund_raising.go/(repo FuncRaisingRepo) InsertDB()/dao.InsertDB()", err.Error())
			}
		}

	}
}

func (dao FuncRaising) InsertDB() error {
	query :=
		`
		INSERT INTO fund_raising
			(projectcode, projectname, projectlogo,
				investorcode, investorname, investorlogo,
				fundstagecode, fundstagename, fundamount,
				funddate, createddate, updateddate,
				description, announcementurl, valulation,
				srcFund, srcInvestor)
		VALUES
			($1, $2, $3,
				$4, $5, $6,
				$7, $8, $9,
				$10, $11, $12,
				$13, $14, $15,
				$16, $17);
	`

	//default
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	_, err := db.PSQL.Exec(query,
		dao.ProjectCode, dao.ProjectName, dao.ProjectLogo,
		dao.InvestorCode, dao.InvestorName, dao.InvestorLogo,
		dao.FundStageCode, dao.FundStageName, dao.FundAmount,
		dao.FundDate, dao.CreatedDate, dao.UpdatedDate,
		dao.Description, dao.AnnouncementUrl, dao.Valulation,
		dao.SrcFund, dao.SrcInvestor)

	return err
}
