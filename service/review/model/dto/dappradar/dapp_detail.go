package dto_dappradar

type DetailDapp struct {
	//From list page
	Image       string
	ProductId   string
	ProductName string

	//
	BlockchainId string
	CategoryId   string
	Tags         []string
	Description  string
	Social       map[string]any
}
