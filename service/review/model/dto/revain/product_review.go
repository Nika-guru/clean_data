package dto

type ProductReviewRepo struct {
	ProductReviews []*ProductReview
}

type ProductReview struct {
	Endpoint  string
	ProductId uint64
}
