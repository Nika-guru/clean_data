package review

import (
	"review-service/pkg/router"
	controller_coingecko "review-service/service/review/controller/coingecko"
	controller_revain "review-service/service/review/controller/revain"

	"github.com/go-chi/chi"
)

var ExchangeInfoServiceSubRoute = chi.NewRouter()

func init() {
	ExchangeInfoServiceSubRoute.Group(func(r chi.Router) {
		r.Get(router.RouterBasePath+`/debug/revain/product-info`, controller_revain.GetProductInfoDebug)
		r.Get(router.RouterBasePath+`/debug/coingecko/product-category`, controller_coingecko.GetProductCategoryDebug)

		/////////////////////////////////////////////////
		r.Get(router.RouterBasePath+"/product-type/all", controller_revain.GetProductTypes)
		r.Get(router.RouterBasePath+"/product-info/search", controller_revain.SearchProductInfoByKeywordAndType)
		r.Get(router.RouterBasePath+"/product-info", controller_revain.GetProductTypes)
	})
}
