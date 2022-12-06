package crawl_coingecko_category

const (
	_baseUrl                   = `https://www.coingecko.com`
	_paramsProductIdByCategory = `?page=%d`
)

func Crawl() {
	CrawlProductCategories()
}
