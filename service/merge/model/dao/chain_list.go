package dao

import "base/pkg/db"

type ChainList struct {
	ChainId   string
	ChainName string
}

func (dao *ChainList) SelectChainNameByChainId() error {
	query := `
	SELECT
		chainname
	FROM 
		chain_list
	WHERE chainid = $1
	`

	rows, err := db.PSQL.Query(query, dao.ChainId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&dao.ChainName)
	}
	return nil
}

func (dao *ChainList) SelectChainIdByChainName() error {
	query := `
	SELECT
		chainid 
	FROM 
		chain_list
	WHERE chainname = $1
	`

	rows, err := db.PSQL.Query(query, dao.ChainName)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&dao.ChainId)
	}
	return nil
}
