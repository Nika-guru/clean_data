package dao

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/db"
	"review-service/pkg/log"
	"review-service/pkg/utils"
)

type BlockchainRepo struct {
	Blockchains []*Blockchain
}

type Blockchain struct {
	Id             uint64
	BlockchainId   string
	BlockchainName string
	Info           map[string]any
	CreatedDate    string
	UpdatedDate    string
}

func (repo *BlockchainRepo) InsertDB() {
	for _, dao := range repo.Blockchains {
		isExist, err := dao.SelectByBlockchainId()
		if err != nil {
			log.Println(log.LogLevelDebug, `service/reivew/model/dao/blockchain.go/func (repo *BlockchainRepo) InsertDB()/ dao.SelectByBlockchainId()`, err.Error())
			continue
		}
		if !isExist {
			err := dao.InsertDB()
			if err != nil {
				log.Println(log.LogLevelError, `service/review/model/dao/blockchain.go/func (repo *BlockchainRepo) InsertDB()/ dao.InsertDB()`, err.Error())
				continue
			}
		}
	}
}

func (dao *Blockchain) SelectByBlockchainId() (isExist bool, err error) {
	query := `
	SELECT DISTINCT ON(blockchainId) 
			id, blockchainId, blockchainname,
			info, createddate, updateddate
		FROM blockchain
		WHERE blockchainId = $1
		ORDER BY blockchainId, id; --get smallest id if duplicate
	`
	rows, err := db.PSQL.Query(query, dao.BlockchainId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		blockchainInfoJSONB := []byte{}
		err := rows.Scan(&dao.Id, &dao.BlockchainId, &dao.BlockchainName,
			&blockchainInfoJSONB, &dao.CreatedDate, &dao.UpdatedDate)
		if err != nil {
			return true, err
		}
		json.Unmarshal(blockchainInfoJSONB, &dao.Info)
		if err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil
}

func (dao *Blockchain) InsertDB() error {
	query := `
	INSERT INTO blockchain
		(blockchainId, blockchainname, info,
			createddate, updateddate)
	VALUES
		($1, $2, $3,
			$4, $5)
	RETURNING id;
	`

	//Set default value
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	blockchainInfoJSONB, err := json.Marshal(dao.Info)
	if err != nil {
		return err
	}

	var id uint64
	err = db.PSQL.QueryRow(query,
		dao.BlockchainId, dao.BlockchainName, blockchainInfoJSONB,
		dao.CreatedDate, dao.UpdatedDate).
		Scan(&id)
	dao.Id = id
	fmt.Println(dao.Id)
	return err

}
