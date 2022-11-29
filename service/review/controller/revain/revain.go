package controller_revain

import (
	"fmt"
	"net/http"
	"review-service/pkg/router"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	dto_revain "review-service/service/review/model/dto/revain"
	"strconv"
	"strings"
)

func GetProductInfoDebug(w http.ResponseWriter, r *http.Request) {
	endpointLength := len(constant.ENDPOINTS_PRODUCT_INFO_REVAIN)
	respMsg := fmt.Sprintf(`Get info coninfo-debug successfully (%d endpoint coin info)`, endpointLength)
	debug := dto_revain.Debug{}
	router.ResponseSuccessWithData(w, "200", respMsg, debug.GetProductInfo())
}

//////////////////////////////////////////////////////////

type DefiInfo struct {
	Image       string `json:"image"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetProductTypes(w http.ResponseWriter, r *http.Request) {
	repo := dao.ProductTypeRepo{}
	err := repo.SelectAll()
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	router.ResponseSuccessWithData(w, "200", "Get info coin successfully", repo.ProductTypes)
}

func SearchProductInfoByKeywordAndType(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	keyword = strings.TrimSpace(keyword)

	typeId := r.URL.Query().Get("typeId")
	if typeId == "" {
		errMsg := "typeId required"
		router.ResponseBadRequest(w, "4xx", errMsg)
		return
	}
	typeIdInt, err := strconv.ParseInt(typeId, 10, 64)
	if err != nil {
		errMsg := "typeId must be integer"
		router.ResponseBadRequest(w, "4xx", errMsg)
		return
	}

	repo := dao.ProductRepo{}
	err = repo.SelectByTitleWithShortDescriptionAndType(keyword, typeIdInt)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	router.ResponseSuccessWithData(w, "200", "Get info coin successfully", repo.Products)
}

func GetProductInfoBySortPagination(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")     //averagestar, comment
	sortType := r.URL.Query().Get("sortType") //asc ^ desc ^ default(desc)
	skip := r.URL.Query().Get("skip")
	limit := r.URL.Query().Get("limit")

	// repo := dao.ProductInfoRepo{}
	// repo.Select()

	fmt.Println(sortBy, sortType, skip, limit)
}
