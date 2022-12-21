package controller

import (
	"crawler/pkg/router"
	"crawler/pkg/utils"
	"crawler/service/merge/model/dao"
	"fmt"
	"net/http"
	"time"
)

type ResponseInfoCrawl struct {
	Source                    string `json:"source"`
	TotalCrawledProduct       uint64 `json:"total_crawled_product"`
	TodayCrawledProduct       uint64 `json:"today_crawled_product"`
	TodayCrawledFailedProduct uint64 `json:"today_crawled_failed_product"`
	LatestCrawledTime         string `json:"latest_crawled_time"`
}

func Info(w http.ResponseWriter, r *http.Request) {
	source := `icoholder`

	product := &dao.Product{}
	rs, err := product.SelectTotal(source)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
	}

	product1 := &dao.Product{}
	rs1, err := product1.SelectTotalIn24h(source)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
	}

	product2 := &dao.Product{}
	rs2, err := product2.SelectLatestCreatedDate(source)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
	}
	latestTime, _ := utils.StringToTimestamp(rs2)
	timeDiff := time.Since(latestTime)

	product3 := &dao.Product{}
	rs3, err := product3.SelectCrawledFailedTotal(source)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
	}

	res := ResponseInfoCrawl{
		Source:                    source,
		TotalCrawledProduct:       rs,
		TodayCrawledProduct:       rs1,
		TodayCrawledFailedProduct: rs3,
		LatestCrawledTime:         fmt.Sprintf("%v ago", timeDiff),
	}

	router.ResponseSuccessWithData(w, "200", "success", res)
}
