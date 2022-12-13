package dto_dappradar

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/db"
	"review-service/pkg/utils"
	"review-service/service/constant"
)

type EndpointDappRepo struct {
	EndpointDappList []EndpointDapp
}

type EndpointDapp struct {
	Endpoint       string
	BlockchainName string

	DetailDapp  *DetailDapp
	CreatedDate string
	UpdatedDate string
}

func (endpointDapp *EndpointDapp) InsertDB() error {
	query :=
		`
		INSERT INTO project 
			(id, "name", category,
			subCategory, social, image, 
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

	var socalsJSONB any //default nil
	var err error
	if endpointDapp.DetailDapp.Social != nil {
		socalsJSONB, err = json.Marshal(endpointDapp.DetailDapp.Social)
		if err != nil {
			return err
		}
	}

	endpointDapp.CreatedDate = utils.Timestamp()
	endpointDapp.UpdatedDate = utils.Timestamp()

	subCategories := ``
	for index, category := range endpointDapp.DetailDapp.SubCategories {
		subCategories += category
		if index == len(endpointDapp.DetailDapp.SubCategories)-1 {
			continue
		}
		subCategories += `,`
	}
	_, err = db.PSQL.Exec(query,
		endpointDapp.DetailDapp.ProductId, endpointDapp.DetailDapp.ProductName, endpointDapp.DetailDapp.Category,
		subCategories, socalsJSONB, endpointDapp.DetailDapp.Image,
		endpointDapp.DetailDapp.Description, `Undefined`, endpointDapp.DetailDapp.ChainName,
		nil, fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, endpointDapp.Endpoint), endpointDapp.CreatedDate,
		endpointDapp.UpdatedDate)
	if err != nil {
		return err
	}
	return nil
}
