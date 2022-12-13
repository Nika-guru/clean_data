package crawl_coincarp_investor

import (
	"encoding/json"
	"io"
	"net/http"
	"review-service/pkg/log"
	"review-service/service/constant"
	"time"
)

const (
	_apiKey = ``
)

type InvestorCoincarpRespCheck struct {
	InvestorData A      `json:"data"`
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
}

type A struct {
	TotalCount int `json:"total_count"`
	TotalPage  int `json:"total_pages"`
	Page       int `json:"page"`
}

type InvestorCoincarpResp struct {
	InvestorData InvestorData `json:"data"`
}
type InvestorData struct {
	InvestorItems []InvestorItem `json:"list"`
}

type InvestorItem struct {
	InvestorCode  string `json:"investorcode"`
	InvestorName  string `json:"investorname"`
	Logo          string `json:"logo"`
	CategoryName  string `json:"categoryname"`
	Location      string `json:"location"`
	Launched      int    `json:"launched"`
	SiteLink      string `json:"sitelink"`
	TwitterLink   string `json:"twitterlink"`
	Linkedin      string `json:"linkedin"`
	EmailAddress  string `json:"emailaddress"`
	FundstageCode string `json:"fundstagecode"`
	InvestorCount string `json:"investorcount"`
}

func (dto *InvestorCoincarpResp) CrawlInvestor(url string) (finish bool, err error) {
beginCallAPI:
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", `Bearer `+_apiKey)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(log.LogLevelDebug, `service/review/crawl/coincarp_investor/crawl_investors.go/CrawlInvestor/client.Do(req)`, err.Error())
		time.Sleep(constant.WAIT_DURATION_WHEN_RATE_LIMIT)
		goto beginCallAPI
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	check := &InvestorCoincarpRespCheck{}
	err = json.Unmarshal(body, &check)
	if err != nil {
		return false, err
	}
	if check.InvestorData.Page > check.InvestorData.TotalPage {
		return true, nil
	}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return false, err
	}

	return false, nil
}
