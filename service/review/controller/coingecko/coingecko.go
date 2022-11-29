package controller_coingecko

import (
	"fmt"
	"net/http"
	"review-service/pkg/router"
	"review-service/service/constant"
	dto_coingecko "review-service/service/review/model/dto/coingecko"
)

func GetProductCategoryDebug(w http.ResponseWriter, r *http.Request) {
	endpointLength := len(constant.ENDPOINTS_PRODUCT_INFO_REVAIN)
	respMsg := fmt.Sprintf(`Get info coninfo-debug successfully (%d endpoint coin info)`, endpointLength)
	debug := dto_coingecko.Debug{}
	router.ResponseSuccessWithData(w, "200", respMsg, debug.GetProductCategory())
}
