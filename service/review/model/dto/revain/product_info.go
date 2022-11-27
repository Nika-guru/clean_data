package dto

type ProductInfoRepo struct {
	Products []*ProductInfo
}

type ProductInfo struct {
	//For product info
	ProductId          *uint64        `json:"product_id"` //point to dao
	ProductName        string         `json:"title"`
	ProductImage       string         `json:"image"`
	ProductDescription string         `json:"description"`
	ProductDetail      map[string]any `json:"detail"`
	CrawlSource        string         `json:"crawl_source"`
	CreatedDate        string         `json:"created_date"`
	UpdatedDate        string         `json:"updated_date"`

	//Add extra for crawl product detail
	EndpointProductDetail string
	//Add extra for category, product category, sub product category
	ProductCategories []string
}
