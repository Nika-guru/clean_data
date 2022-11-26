package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/log"
)

type Blockchain struct {
	ChainId     string `json:"chainId"`
	ChainSymbol string `json:"chainSymbol"`
	ChainName   string `json:"chainName"`
}

func (blockchain *Blockchain) SelectBlockchainByChainname() error {
	query :=
		`
		SELECT chain_id, chain_symbol, chain_name
		FROM blockchains WHERE chain_name= $1
		`
	rows, err := db.PSQL.Query(query, blockchain.ChainName)
	if err != nil {
		log.Println(log.LogLevelError, "dao/blockchain/SelectBlockchainByChainname", "Query db error")
		return err
	}

	defer rows.Close()
	if rows.Next() {
		//NOT NULL in db, data fixed then don't need
		rows.Scan(&blockchain.ChainId, &blockchain.ChainSymbol, &blockchain.ChainName)
	}

	return nil
}
