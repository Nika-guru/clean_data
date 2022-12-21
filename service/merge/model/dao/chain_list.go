package dao

import "crawler/pkg/db"

type ChainListRepo struct {
	ChainLists []ChainList
}
type ChainList struct {
	ChainId   string
	ChainName string
}

func (repo *ChainListRepo) SelectAll() error {
	query :=
		`
	SELECT chainid, chainname
	FROM chain_list;	
	`
	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		dao := ChainList{}
		rows.Scan(&dao.ChainId, &dao.ChainName)
		repo.ChainLists = append(repo.ChainLists, dao)
	}

	return nil
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
