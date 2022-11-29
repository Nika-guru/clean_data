package crawl_coingecko

const (
	_baseUrl                   = `https://www.coingecko.com`
	_paramsProductIdByCategory = `?page=%d`
)

func init() {
	go CrawlProductCategories()
}
