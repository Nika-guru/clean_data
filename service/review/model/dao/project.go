package dao

import (
	"review-service/pkg/db"
	"review-service/pkg/utils"
)

type Project struct {
	ProjectCode     string
	ProjectName     string
	ProjectCategory string
	Project         string
	Src             string
	CreatedDate     string
	UpdatedDate     string
}

func (dao *Project) InsertDB() error {
	query :=
		`
	INSERT INTO project
		(id, "name", category, 
		subcategory, social, image, 
		description, chainid, chainname, 
		extradata, src, createddate, 
		updateddate)
	VALUES
		($1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			$13);		
`

	//default
	dao.CreatedDate = utils.Timestamp()
	dao.UpdatedDate = utils.Timestamp()

	_, err := db.PSQL.Query(query)

	return err
}
