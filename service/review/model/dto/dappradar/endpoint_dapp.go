package dto_dappradar

type EndpointDappRepo struct {
	EndpointDappList []EndpointDapp
}

type EndpointDapp struct {
	Endpoint       string
	BlockchainName string

	DetailDapp *DetailDapp
}
