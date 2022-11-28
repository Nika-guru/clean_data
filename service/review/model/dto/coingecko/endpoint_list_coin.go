package dto_coingecko

import "review-service/service/review/model/dao"

type EndpointCategoryRepo struct {
	EndpointCategories []*EndpointCategory
}

type EndpointCategory struct {
	CategoryName string
	Endpoint     string
	CoinIdList   []string
}

func (dto *EndpointCategory) ConvertTo(dao.Category) {

}
