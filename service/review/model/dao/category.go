package dao

import (
	"fmt"
	"review-service/pkg/db"
)

type Category struct {
	CategoryId   uint64
	CategoryName string
}

func (category *Category) SelectByName() (isExist bool, err error) {
	query := `
		SELECT DISTINCT ON("name") id, "name"
		FROM category
		WHERE "name" = $1
		ORDER BY "name", id; --get smallest id if duplicate
	`
	rows, err := db.PSQL.Query(query, category.CategoryName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, rows.Scan(&category.CategoryId, &category.CategoryName)
	}

	return false, nil
}

func (category *Category) InsertDB() error {
	query := `
	INSERT INTO category
		("name")
	VALUES
		($1)

		RETURNING id;
	`
	fmt.Println(`====run here====`, category.CategoryName)

	var categoryId uint64
	err := db.PSQL.QueryRow(query,
		category.CategoryName).Scan(&categoryId)
	category.CategoryId = categoryId
	return err
}
