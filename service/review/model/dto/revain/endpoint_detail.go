package dto

type EndpointDetailRepo struct {
	Endpoints []EndpointDetail
}

type EndpointDetail struct {
	ProductId uint64
	Endpoint  string
}
