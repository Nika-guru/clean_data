package dao

import (
	"database/sql"
	"reflect"
	"review-service/pkg/db"
)

type ContractsInfo struct {
	CoinId    string
	Contracts []*Contract
}

type Contract struct {
	ChainName    *string
	TokenAddress *string
	ChainImage   *string
}

type ContractDB struct {
	ChainName    sql.NullString
	TokenAddress sql.NullString
	ChainImage   sql.NullString
}

func (contractDB *ContractDB) convertTo(contract *Contract) {
	if reflect.TypeOf(contractDB.ChainName) != nil {
		//######### ChainName is not null value#########
		if contractDB.ChainName.Valid {
			contract.ChainName = &contractDB.ChainName.String
		}
	}

	if reflect.TypeOf(contractDB.TokenAddress) != nil {
		//######### TokenAddress is not null value#########
		if contractDB.TokenAddress.Valid {
			contract.TokenAddress = &contractDB.TokenAddress.String
		}
	}

	if reflect.TypeOf(contractDB.ChainImage) != nil {
		//######### ChainImage is not null value#########
		if contractDB.ChainImage.Valid {
			contract.ChainImage = &contractDB.ChainImage.String
		}
	}
}

func (contractsInfo *ContractsInfo) SelectAllContractByCoinId() error {
	tx, err := db.PSQL.Begin()
	if err != nil {
		return err
	}
	query :=
		`
		SELECT token_coins.chain_name, token_coins.token_address, native_coins.coin_image as chain_image
		FROM coins AS token_coins
			LEFT join coins AS native_coins
			ON token_coins.chain_name = native_coins.coin_id
			WHERE (token_coins.chain_name is not null and token_coins.token_address is not null ) --token coin only
			AND token_coins.coin_id = $1
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(contractsInfo.CoinId)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		contractDB := &ContractDB{}
		rows.Scan(&contractDB.ChainName, &contractDB.TokenAddress, &contractDB.ChainImage)
		contract := &Contract{}
		contractDB.convertTo(contract)

		contractsInfo.Contracts = append(contractsInfo.Contracts, contract)
	}

	return tx.Commit()
}
