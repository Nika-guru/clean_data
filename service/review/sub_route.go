package review

import (
	"review-service/pkg/router"
	"review-service/service/review/controller"

	"github.com/go-chi/chi"
)

var ExchangeInfoServiceSubRoute = chi.NewRouter()

func init() {
	ExchangeInfoServiceSubRoute.Group(func(r chi.Router) {
		r.Get(router.RouterBasePath+"/product-type/all", controller.GetProductTypes)
		r.Get(router.RouterBasePath+"/product-info/search", controller.SearchProductInfoByKeywordAndType)
		r.Get(router.RouterBasePath+"/product-info", controller.GetProductTypes)
	})
}
