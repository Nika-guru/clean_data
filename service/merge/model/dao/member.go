package dao

import (
	"crawler/pkg/db"
	"crawler/pkg/log"
	"crawler/pkg/utils"
	"database/sql"
	"encoding/json"
)

type MemberRepo struct {
	Members []*Member
}

type Member struct {
	MemberName string
	Detail     map[string]any //	MemberImage    string;	MemberLinkin   string;	MemberPosition string;	VerifiedMember string, ProductCodeIcoHolder

	ProductId     uint64
	ProductName   string
	ProductSymbol string
	CreatedDate   string
	UpdatedDate   string
}

func (repo *MemberRepo) InsertDB() {
	for _, member := range repo.Members {
		isExist, err := member.IsExist()
		if err != nil {
			log.Println(log.LogLevelError, `service/merge/model/dao/member.go/ (repo *MemberRepo) InsertDB()/ member.IsExist(), at member name `+member.MemberName, err.Error())
			continue
		}

		if !isExist {
			err = member.InsertDB()
			if err != nil {
				log.Println(log.LogLevelError, `service/merge/model/dao/member.go/ (repo *MemberRepo) InsertDB()/ member.InsertDB(), at member name  `+member.MemberName, err.Error())
				continue
			}
		}

	}
}

func (dao *Member) IsExist() (bool, error) {
	query :=
		`
		SELECT 
			*
		FROM
			member
		WHERE 
		
		`
	var rows *sql.Rows
	var err error
	if dao.Detail[`memberLinkedin`] != nil {
		query += ` (detail->>'memberLinkedin' = $1 and detail->>'src' = $2) `
		rows, err = db.PSQL.Query(query, dao.Detail[`memberLinkedin`], dao.Detail[`src`])
	} else
	//Certain insert when have no linkedin to check
	{
		query += ` 1 != 1 `
		rows, err = db.PSQL.Query(query)
	}
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (dao *Member) InsertDB() error {
	query := `
	INSERT INTO "member"
		(membername, productid, productname,
			productsymbol, createddate, updateddate,
			detail)
	VALUES
		($1, $2, $3,
			$4, $5, $6,
			$7);
	`
	//Default value
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	detailJSON, err := json.Marshal(dao.Detail)
	if err != nil {
		return err
	}

	_, err = db.PSQL.Exec(query,
		dao.MemberName, dao.ProductId, dao.ProductName,
		dao.ProductSymbol, dao.CreatedDate, dao.UpdatedDate,
		detailJSON)
	return err
}
