package dao

import (
	"database/sql"
	"review-service/pkg/db"
	"review-service/pkg/log"
)

type LinkGroupInfo struct {
	CoinId    string
	LinkItems map[string][]*LinkItem
}

type LinkItem struct {
	LinkTitle *string `json:"link_title"`
	LinkIcon  *string `json:"link_icon"`
	LinkUrl   *string `json:"link_url"`
}

type RowDB struct {
	GroupName sql.NullString
	LinkTitle sql.NullString
	LinkIcon  sql.NullString
	LinkUrl   sql.NullString
}

func (linkGroup *LinkGroupInfo) SelectLinkGroupsByCoinId() error {
	query :=
		`
		SELECT link_groups.name as group_name, link_items.link_title, link_items.link_icon, 
				link_items.link_url  from link_items 
			INNER JOIN link_groups 
			ON link_items.link_group_id = link_groups.id
		WHERE link_items.coin_id = $1 --detail coin
		ORDER BY link_groups.id, link_items.id --query by order insert FIFO
	`

	rows, err := db.PSQL.Query(query, linkGroup.CoinId)
	if err != nil {
		errMsg := "query database error"
		log.Println(log.LogLevelWarn, "CoinDetail/SelectDetailInfoByCoinId, detail: "+err.Error(), errMsg)
		return err
	}
	defer rows.Close()

	linkGroup.LinkItems = map[string][]*LinkItem{}
	for rows.Next() {
		rowDB := RowDB{}
		rows.Scan(&rowDB.GroupName, &rowDB.LinkTitle, &rowDB.LinkIcon,
			&rowDB.LinkUrl)

		//Data not nil
		if rowDB.GroupName.Valid {
			groupItems, found := linkGroup.LinkItems[rowDB.GroupName.String]

			if !found {
				groupItems = []*LinkItem{}
			}

			groupItem := &LinkItem{}
			if rowDB.LinkTitle.Valid {
				groupItem.LinkTitle = &rowDB.LinkTitle.String
			}
			if rowDB.LinkIcon.Valid {
				groupItem.LinkIcon = &rowDB.LinkIcon.String
			}
			if rowDB.LinkUrl.Valid {
				groupItem.LinkUrl = &rowDB.LinkUrl.String
			}
			groupItems = append(groupItems, groupItem)
			linkGroup.LinkItems[rowDB.GroupName.String] = groupItems
		}
	}

	return nil
}
