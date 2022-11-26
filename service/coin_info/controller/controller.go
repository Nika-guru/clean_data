package controller

import (
	"net/http"
	"review-service/pkg/log"
	"review-service/pkg/router"
	daoView "review-service/service/coin_info/model/dao/view"

	"github.com/go-chi/chi"
)

func GetCoinInfo(w http.ResponseWriter, r *http.Request) {
	//Path url type only is string
	coinId := chi.URLParam(r, "coinId")
	if coinId == "" {
		errMsg := "coinId required"
		log.Println(log.LogLevelDebug, "GetLatestCoinPrice: URLParam(r, \"coinId\")", errMsg)
		router.ResponseBadRequest(w, "4xx", errMsg)
		return
	}

	response := make(map[string]interface{})

	detailInfo := &daoView.DetailInfo{}
	detailInfo.CoinId = &coinId
	detailInfo.SelectDetailByCoinId()
	response["coin"] = detailInfo

	linkGroup := &daoView.LinkGroupInfo{}
	linkGroup.CoinId = coinId
	linkGroup.SelectLinkGroupsByCoinId()
	response["links"] = linkGroup.LinkItems

	contractsInfo := &daoView.ContractsInfo{}
	contractsInfo.CoinId = coinId
	contractsInfo.SelectAllContractByCoinId()
	response["contracts"] = linkGroup.LinkItems

	tagInfo := &daoView.TagInfo{}
	tagInfo.CoinId = coinId
	tagInfo.SelectTagsByCoinId()
	response["tags"] = tagInfo.Tags

	router.ResponseSuccessWithData(w, "200", "Get info coin successfully", response)
}
