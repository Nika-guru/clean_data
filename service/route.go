package service

import (
	"review-service/pkg/router"
	"review-service/service/index"
	"review-service/service/review"

	//1
	// crawl_coingecko_category "review-service/service/review/crawl/coingecko_category"
	// crawl_revain_category_productinfo_review_reply "review-service/service/review/crawl/revain_category_productinfo_review_reply"

	//2
	crawl_dappradar_dapp "review-service/service/review/crawl/dappradar_dapp"
)

func init() {
	go func() {
		//1
		// crawl_revain_category_productinfo_review_reply.Crawl()
		// crawl_coingecko_category.Crawl()

		//2
		crawl_dappradar_dapp.Crawl()
	}()
}

// LoadRoutes to Load Routes to Router
func LoadRoutes() {

	// Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	router.Router.Mount(router.RouterBasePath+"/review", review.ExchangeInfoServiceSubRoute)
}
