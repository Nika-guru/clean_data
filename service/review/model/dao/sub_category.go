package dao

import (
	"review-service/pkg/db"
)

type SubCategory struct {
	SubCategoryId   uint64
	SubCategoryName string
	CategoryId      uint64
}

func (subCategory *SubCategory) SelectByName() (isExist bool, err error) {
	query := `
		SELECT DISTINCT ON("name") id, "name", categoryid
		FROM sub_category
		WHERE "name" = $1
		ORDER BY "name", id;
	`
	rows, err := db.PSQL.Query(query, subCategory.SubCategoryName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, rows.Scan(&subCategory.SubCategoryId, &subCategory.SubCategoryName, &subCategory.CategoryId)
	}

	return false, nil
}

func (subCategory *SubCategory) InsertDB() error {
	query := `
	INSERT INTO sub_category
		("name", categoryid)
	VALUES
		($1, $2)
	
		RETURNING id;
	`

	var subCategoryId uint64
	err := db.PSQL.QueryRow(query,
		subCategory.SubCategoryName, subCategory.CategoryId).Scan(&subCategoryId)
	subCategory.SubCategoryId = subCategoryId
	return err
}
