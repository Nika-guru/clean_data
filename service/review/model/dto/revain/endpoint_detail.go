package dto_revain

type EndpointDetailRepo struct {
	Endpoints []EndpointDetail
}

type EndpointDetail struct {
	ProductId uint64
	Endpoint  string
}

func (dtoEndpointRepo *EndpointDetailRepo) ConvertFrom(dtoProductInfoRepo *ProductInfoRepo) {
	for _, dtoProduct := range dtoProductInfoRepo.Products {
		dtoEndpointRepo.Endpoints = append(dtoEndpointRepo.Endpoints, EndpointDetail{
			ProductId: *dtoProduct.ProductId,
			Endpoint:  dtoProduct.EndpointProductDetail,
		})
	}
}
