package dao

import (
	"database/sql"
	"reflect"
	"review-service/pkg/db"
)

type DetailInfo struct {
	CoinId          *string
	CoinSymbol      *string
	CoinName        *string
	CoinImage       *string
	CoinDescription *string
}

type DetailInfoDB struct {
	CoinId          sql.NullString
	CoinSymbol      sql.NullString
	CoinName        sql.NullString
	CoinImage       sql.NullString
	CoinDescription sql.NullString
}

func (detailInfoDB *DetailInfoDB) convertTo(cetailInfo *DetailInfo) {
	if reflect.TypeOf(detailInfoDB.CoinId) != nil {
		//######### CoinId is not null value#########
		if detailInfoDB.CoinId.Valid {
			cetailInfo.CoinId = &detailInfoDB.CoinId.String
		}
	}

	if reflect.TypeOf(detailInfoDB.CoinSymbol) != nil {
		//######### CoinSymbol is not null value#########
		if detailInfoDB.CoinSymbol.Valid {
			cetailInfo.CoinSymbol = &detailInfoDB.CoinSymbol.String
		}
	}

	if reflect.TypeOf(detailInfoDB.CoinName) != nil {
		//######### CoinName is not null value#########
		if detailInfoDB.CoinName.Valid {
			cetailInfo.CoinName = &detailInfoDB.CoinName.String
		}
	}

	if reflect.TypeOf(detailInfoDB.CoinImage) != nil {
		//######### CoinImage is not null value#########
		if detailInfoDB.CoinImage.Valid {
			cetailInfo.CoinImage = &detailInfoDB.CoinImage.String
		}
	}

	if reflect.TypeOf(detailInfoDB.CoinDescription) != nil {
		//######### CoinDescription is not null value#########
		if detailInfoDB.CoinDescription.Valid {
			cetailInfo.CoinDescription = &detailInfoDB.CoinDescription.String
		}
	}
}
func (detailInfo *DetailInfo) SelectDetailByCoinId() error {
	tx, err := db.PSQL.Begin()
	if err != nil {
		return err
	}
	query :=
		`
		select distinct on(coin_id) coin_id, coin_symbol, coin_name, coin_image, coin_description from coins
		where coin_id = $1
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(detailInfo.CoinId)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		detailInfoDB := &DetailInfoDB{}
		rows.Scan(&detailInfoDB.CoinId, &detailInfoDB.CoinSymbol, &detailInfoDB.CoinName,
			&detailInfoDB.CoinImage, &detailInfoDB.CoinDescription)
		detailInfoDB.convertTo(detailInfo)
	}

	return tx.Commit()
}
