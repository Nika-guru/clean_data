package dto_dappradar

type DetailDapp struct {
	//From list page
	Image       string
	ProductId   string
	ProductName string
	//From url
	ChainName string
	Category  string

	//Detail
	SubCategories []string //tags
	Description   string
	Social        map[string]any //url, image
}
