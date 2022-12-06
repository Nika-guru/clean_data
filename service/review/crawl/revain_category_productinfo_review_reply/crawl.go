package crawl_revain_category_productinfo_review_reply

const (
	_baseUrl                      = `https://revain.org`
	_urlQueryParamsProductsInfo   = `?sortBy=reviews&page=%d`
	_urlQueryParamsProductReviews = `?page=%d&sortBy=recent&direction=ASC`
)

func Crawl() {
	CrawlProductsInfo()
}
