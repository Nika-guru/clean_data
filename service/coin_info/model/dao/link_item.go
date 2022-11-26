package dao

import (
	"database/sql"
	"fmt"
	"reflect"
	"review-service/pkg/db"
	"review-service/pkg/log"
)

type LinkItemRepo struct {
	LinkItems []*LinkItem
}
type LinkItem struct {
	LinkTitle   *string
	LinkUrl     *string
	LinkIcon    *string
	LinkGroupId *int
	CoinId      *string
}

type LinkItemDB struct {
	LinkTitle   sql.NullString
	LinkUrl     sql.NullString
	LinkIcon    sql.NullString
	LinkGroupId sql.NullInt32
	CoinId      sql.NullString
}

func (linkItem *LinkItem) convertTo(linkItemDB *LinkItemDB) {
	//######### LinkTitle is not null value  #########
	if linkItem.LinkTitle != nil {
		linkItemDB.LinkTitle = sql.NullString{String: *linkItem.LinkTitle, Valid: true}
	} else {
		linkItemDB.LinkTitle = sql.NullString{Valid: false}
	}

	//######### LinkUrl is not null value  #########
	if linkItem.LinkUrl != nil {
		linkItemDB.LinkUrl = sql.NullString{String: *linkItem.LinkUrl, Valid: true}
	} else {
		linkItemDB.LinkUrl = sql.NullString{Valid: false}
	}

	//######### LinkIcon is not null value  #########
	if linkItem.LinkIcon != nil {
		linkItemDB.LinkIcon = sql.NullString{String: *linkItem.LinkIcon, Valid: true}
	} else {
		linkItemDB.LinkIcon = sql.NullString{Valid: false}
	}

	//######### LinkGroupId is not null value  #########
	if linkItem.LinkGroupId != nil {
		linkItemDB.LinkGroupId = sql.NullInt32{Int32: int32(*linkItem.LinkGroupId), Valid: true}
	} else {
		linkItemDB.LinkGroupId = sql.NullInt32{Valid: false}
	}

	//######### CoinId is not null value  #########
	if linkItem.CoinId != nil {
		linkItemDB.CoinId = sql.NullString{String: *linkItem.CoinId, Valid: true}
	} else {
		linkItemDB.CoinId = sql.NullString{Valid: false}
	}
}

func (linkItemDB *LinkItemDB) convertTo(linkItem *LinkItem) {
	if reflect.TypeOf(linkItemDB.LinkTitle) != nil {
		//######### LinkTitle is not null value#########
		if linkItemDB.LinkTitle.Valid {
			linkItem.LinkTitle = &linkItemDB.LinkTitle.String
		}
	}

	if reflect.TypeOf(linkItemDB.LinkUrl) != nil {
		//######### LinkUrl is not null value#########
		if linkItemDB.LinkUrl.Valid {
			linkItem.LinkUrl = &linkItemDB.LinkUrl.String
		}
	}

	if reflect.TypeOf(linkItemDB.LinkIcon) != nil {
		//######### LinkIcon is not null value#########
		if linkItemDB.LinkIcon.Valid {
			linkItem.LinkIcon = &linkItemDB.LinkIcon.String
		}
	}

	if reflect.TypeOf(linkItemDB.LinkGroupId) != nil {
		//######### LinkGroupId is not null value#########
		if linkItemDB.LinkGroupId.Valid {
			val := int(linkItemDB.LinkGroupId.Int32)
			linkItem.LinkGroupId = &val
		}
	}

	if reflect.TypeOf(linkItemDB.CoinId) != nil {
		//######### CoinId is not null value#########
		if linkItemDB.CoinId.Valid {
			linkItem.CoinId = &linkItemDB.CoinId.String
		}
	}

}

func (repo *LinkItemRepo) InsertDB() (insertedTotal int) {
	//########## Set Inserted Total of query = 0 ##########
	insertedTotal = 0

	//########## Traverse each crawled Coin ##########
	for _, linkItem := range repo.LinkItems {

		linkItemDB := &LinkItemDB{}
		linkItem.convertTo(linkItemDB)

		query :=
			`
			INSERT INTO public.link_items
			(	
				link_title, link_url, link_icon, 
				link_group_id, coin_id
			)
			VALUES($1, $2, $3, $4, $5);		
			`
		_, err := db.PSQL.Exec(query,
			linkItemDB.LinkTitle, linkItemDB.LinkUrl, linkItemDB.LinkIcon,
			linkItemDB.LinkGroupId, linkItemDB.CoinId)

		//########## Insert database failed ##########
		if err != nil {
			log.Println(log.LogLevelError, "dao/link_item.go/InsertDB", err.Error())
			continue
		} else
		//########## Insert database successful ##########
		{
			insertedTotal += 1
		}

	}

	return insertedTotal
}

func (linkItem *LinkItem) GetKeyMap() string {
	key := ""

	if linkItem.CoinId != nil {
		key += *linkItem.CoinId
	}

	if linkItem.LinkGroupId != nil {
		key += fmt.Sprintf("%d", *linkItem.LinkGroupId)
	}

	if linkItem.LinkUrl != nil {
		key += *linkItem.LinkUrl
	}

	return key
}

func (repo *LinkItemRepo) SelectAll() {
	query :=
		`
		SELECT link_title, link_url, link_icon, 
				link_group_id, coin_id
		FROM link_items;		
		`
	rows, err := db.PSQL.Query(query)
	if err != nil {
		log.Println(log.LogLevelError, "dao/link_item.go/SelectAll", "Query database error, detail: "+err.Error())
	}
	for rows.Next() {
		linkItemDB := &LinkItemDB{}
		rows.Scan(&linkItemDB.LinkTitle, &linkItemDB.LinkUrl, &linkItemDB.LinkIcon,
			&linkItemDB.LinkGroupId, &linkItemDB.CoinId)
		linkItem := &LinkItem{}
		linkItemDB.convertTo(linkItem)
		repo.LinkItems = append(repo.LinkItems, linkItem)
	}
}
