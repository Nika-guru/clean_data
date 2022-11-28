package dto_revain

type ProductReviewRepo struct {
	ProductReviews []*ProductReview
}

type ProductReview struct {
	Endpoint  string
	ProductId uint64

	Content    string
	Star       float64
	ReviewDate string

	//Account
	AccountImage string
	Username     string
}
