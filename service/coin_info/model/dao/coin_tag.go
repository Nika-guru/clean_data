package dao

import (
	"database/sql"
	"reflect"
	"review-service/pkg/db"
	"review-service/pkg/log"
)

type CoinTagRepo struct {
	CoinsTags []*CoinTag
}

type CoinTag struct {
	CoinId  *string
	TagName *string
}

type CoinTagDB struct {
	CoinId  sql.NullString
	TagName sql.NullString
}

func (coinTag *CoinTag) ConvertTo(coinTagDB *CoinTagDB) {
	//######### CoinId is not null value  #########
	if coinTag.CoinId != nil {
		coinTagDB.CoinId = sql.NullString{String: *coinTag.CoinId, Valid: true}
	} else {
		coinTagDB.CoinId = sql.NullString{Valid: false}
	}

	//######### TagName is not null value  #########
	if coinTag.TagName != nil {
		coinTagDB.TagName = sql.NullString{String: *coinTag.TagName, Valid: true}
	} else {
		coinTagDB.TagName = sql.NullString{Valid: false}
	}
}

func (coinTagDB *CoinTagDB) ConvertTo(coinTag *CoinTag) {
	if reflect.TypeOf(coinTagDB.CoinId) != nil {
		//######### CoinID is not null value#########
		if coinTagDB.CoinId.Valid {
			coinTag.CoinId = &coinTagDB.CoinId.String
		}
	}

	if reflect.TypeOf(coinTagDB.TagName) != nil {
		//######### TagName is not null value#########
		if coinTagDB.TagName.Valid {
			coinTag.TagName = &coinTagDB.TagName.String
		}
	}
}

func (coinTag *CoinTag) GetKeyCoinsTags() string {
	key := ""
	if coinTag == nil {
		return key
	}

	if coinTag.CoinId != nil {
		key += *coinTag.CoinId
	}

	if coinTag.TagName != nil {
		key += *coinTag.TagName
	}
	return key
}

func (coinTagRepo *CoinTagRepo) InsertDB() (insertedTag, insertedCoinTag int) {
	insertedTag = 0
	insertedCoinTag = 0

	//############ Start existed tags #############
	existedTagsRepo := &CoinTagRepo{}
	err := existedTagsRepo.SelectExistedUniqueTags()
	if err != nil {
		log.Println(log.LogLevelInfo, "dao/coin_tag.go", "coinTagRepo InsertDB : Failed")
		return 0, 0
	}
	//Convert map to search
	existedTagsMap := make(map[string]bool)
	for _, coinTagRunner := range existedTagsRepo.CoinsTags {
		_, foundTag := existedTagsMap[*coinTagRunner.TagName]
		if !foundTag {
			existedTagsMap[*coinTagRunner.TagName] = true
		}
	}
	//############ End existed tags #############

	//############ Start existed coins tags #############
	existedCoinsTagsRepo := &CoinTagRepo{}
	err = existedCoinsTagsRepo.SelectExistedCoinsTags()
	if err != nil {
		log.Println(log.LogLevelInfo, "dao/coin_tag.go", "coinTagRepo InsertDB : Failed")
		return 0, 0
	}
	//Convert map to search
	existedCoinsTagsMap := make(map[string]bool)
	for _, coinTagRunner := range existedCoinsTagsRepo.CoinsTags {
		_, foundTag := existedCoinsTagsMap[coinTagRunner.GetKeyCoinsTags()]
		if !foundTag {
			existedCoinsTagsMap[coinTagRunner.GetKeyCoinsTags()] = true
		}
	}
	//############ End existed coins tags #############

	for _, coinTag := range coinTagRepo.CoinsTags {
		if coinTag != nil {
			_, foundTag := existedTagsMap[*coinTag.TagName]
			if !foundTag {
				err := coinTag.InsertTagDB()
				//Insert tags table fail
				if err != nil {
					log.Println(log.LogLevelInfo, "dao/coin_tag.go", "InsertDB Tag : Failed"+err.Error())
					continue
				}
				insertedTag++
			}

			_, foundCoinTag := existedCoinsTagsMap[coinTag.GetKeyCoinsTags()]
			if !foundCoinTag {
				err := coinTag.InsertCoinTagDB()
				//Insert coins_tags table fail
				if err != nil {
					log.Println(log.LogLevelInfo, "dao/coin_tag.go", "InsertDB CoinTag : Failed"+err.Error())
					continue
				}
				insertedCoinTag++
			}
		}
	}

	return insertedTag, insertedCoinTag
}

func (coinTag *CoinTag) InsertTagDB() error {
	tx, err := db.PSQL.Begin()
	var stmt *sql.Stmt
	if err != nil {
		return err
	}

	query :=
		`
			INSERT INTO tags
				("name")
				VALUES($1);

			`
	stmt, err = tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(coinTag.TagName)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (coinTag *CoinTag) InsertCoinTagDB() error {
	tx, err := db.PSQL.Begin()
	var stmt *sql.Stmt
	if err != nil {
		return err
	}

	query :=
		`
		INSERT INTO public.coins_tags
			(coin_id, tag_name)
			VALUES($1, $2);
	`
	stmt, err = tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(coinTag.CoinId, coinTag.TagName)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *CoinTagRepo) SelectExistedUniqueTags() error {
	query :=
		`
		SELECT DISTINCT ON(name) name FROM tags
	`
	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		coinTagDB := &CoinTagDB{}
		rows.Scan(&coinTagDB.TagName)
		coinTag := &CoinTag{}
		coinTagDB.ConvertTo(coinTag)
		repo.CoinsTags = append(repo.CoinsTags, coinTag)
	}

	return nil
}

func (repo *CoinTagRepo) SelectExistedCoinsTags() error {
	tx, err := db.PSQL.Begin()
	if err != nil {
		return err
	}
	query :=
		`
		SELECT coin_id, tag_name
		FROM coins_tags;
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		coinTagDB := &CoinTagDB{}
		rows.Scan(&coinTagDB.CoinId, &coinTagDB.TagName)
		coinTag := &CoinTag{}
		coinTagDB.ConvertTo(coinTag)
		repo.CoinsTags = append(repo.CoinsTags, coinTag)
	}

	return tx.Commit()
}
