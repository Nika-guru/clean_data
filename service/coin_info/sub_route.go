package coin_info

import (
	"review-service/pkg/router"
	"review-service/service/coin_info/controller"

	"github.com/go-chi/chi"
)

var CoinInfoServiceSubRoute = chi.NewRouter()

func init() {
	CoinInfoServiceSubRoute.Group(func(r chi.Router) {
		r.Get(router.RouterBasePath+"/{coinId}", controller.GetCoinInfo)
	})
}
