package dao

import (
	"database/sql"
	"review-service/pkg/db"
)

type TagInfo struct {
	CoinId string
	Tags   []*string
}

type RowTagDB struct {
	TagName sql.NullString
}

func (tagInfo *TagInfo) SelectTagsByCoinId() error {
	tx, err := db.PSQL.Begin()
	if err != nil {
		return err
	}
	query :=
		`
		SELECT  tag_name
			FROM coins_tags
		WHERE coin_id = $1		
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tagInfo.CoinId)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rowTagDB := &RowTagDB{}
		rows.Scan(&rowTagDB.TagName)

		if rowTagDB.TagName.Valid {
			tagInfo.Tags = append(tagInfo.Tags, &rowTagDB.TagName.String)
		}
	}

	return tx.Commit()
}
