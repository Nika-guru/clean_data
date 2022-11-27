package dto

type ProductDetail struct {
	ProductId   uint64
	Description string
	Url         map[string]string //key: url, value: type
	CreatedDate string
	UpdatedDate string
}
