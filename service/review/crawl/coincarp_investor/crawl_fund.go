package crawl_coincarp_investor

import (
	"encoding/json"
	"io"
	"net/http"
	"review-service/pkg/log"
	"review-service/service/constant"
	"time"
)

type FundCoincarpResp struct {
	FundData FundData `json:"data"`
	Code     uint8    `json:"code"`
}

type FundData struct {
	FundItems []FundItem `json:"list"`
}

type FundItem struct {
	ProjectCode   string         `json:"projectcode"`
	ProjectName   string         `json:"projectname"`
	ProjectLogo   string         `json:"logo"`
	CategoryList  []FundCategory `json:"categorylist"`
	EndpointFund  string         `json:"fundcode"`
	FundStageCode string         `json:"fundstagecode"`
	FundStageName string         `json:"fundstagename"`
	FundAmount    float64        `json:"fundamount"`    //Invest
	Valulation    float64        `json:"valulation"`    //Value
	FundDate      uint64         `json:"funddate"`      //Unix
	InvestorCount uint8          `json:"investorcount"` //All investor invest this project in this fund stage
	//investorlist khong lay vi khong du, dung them API

	//Append crawl HTML
	Description     string
	AnnouncementUrl string
	//from url
	SrcFund string
}

type FundCategory struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (dto *FundCoincarpResp) CrawlFund(url string) error {
beginCallAPI:
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", `Bearer `+_apiKey)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(log.LogLevelDebug, `service/review/crawl/coincarp_investor/crawl_fund.go/CrawlFund/client.Do(req)`, err.Error())
		time.Sleep(constant.WAIT_DURATION_WHEN_RATE_LIMIT)
		goto beginCallAPI
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return err
	}

	return nil
}
