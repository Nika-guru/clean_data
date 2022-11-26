package dao

import "review-service/pkg/db"

type ProductTypeRepo struct {
	ProductTypes []ProductType
}

type ProductType struct {
	Id   int    `json:"product_type_id"`
	Name string `json:"product_type_name"`
}

func (repo *ProductTypeRepo) SelectAll() error {
	query := `
		SELECT id, "name"
		FROM product_type;
	`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}

	//Mock data (Default value for all type)
	repo.ProductTypes = append(repo.ProductTypes, ProductType{
		Id:   -1,
		Name: `All`,
	})

	defer rows.Close()
	for rows.Next() {
		productType := ProductType{}
		rows.Scan(&productType.Id, &productType.Name)

		repo.ProductTypes = append(repo.ProductTypes, productType)
	}
	return nil
}
