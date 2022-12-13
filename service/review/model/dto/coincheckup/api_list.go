package dto_coincheckup

//https://coincheckup.com/api/coincheckup/get_coin_list?order_by=last_market_cap_usd&order_direction=desc&limit=10&offset=1&t=27838521
type APIList struct {
	Data []Data `json:"data"`
}

type Data struct {
	CoinId     string `json:"shortname"`
	CoinSymbol string `json:"symbol"`
	CoinName   string `json:"name"`
	Endpoint   string `json:"ccu_slug"`

	Detail *APIDetail
}
