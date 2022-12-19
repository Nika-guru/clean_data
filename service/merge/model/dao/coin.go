package dao

import "strings"

type CoinRepo struct {
	Coins []*Coin
}

type Coin struct {
	CoinCode      string         `json:"coinId"`  //coin id
	Type          string         `json:"type"`    //Coin or Token
	TokenAddress  string         `json:"address"` //lowcase
	ChainId       string         `json:"chainId"`
	ChainName     string         `json:"chainName"`
	CoinSymbol    string         `json:"symbol"` //upcase
	CoinName      string         `json:"name"`
	Tag           string         `json:"tag"`
	Decimals      uint8          `json:"decimals"`
	TotalSupply   string         `json:"totalSupply"`
	MaxSupply     string         `json:"maxSupply"`
	Marketcap     string         `json:"marketcap"`
	VolumeTrading string         `json:"volumeTrading"`
	CoinImage     string         `json:"image"`
	Src           string         `json:"source"`
	Detail        map[string]any `json:"detail"`
	CreatedDate   string         `json:"createdDate"`
	UpdatedDate   string         `json:"updatedDate"`
	Holder        int            //Crawl explorer
	Description   string
	Contract      map[string]any
}

func (daoCoin *Coin) ConvertToProduct(daoProduct *Product) {
	daoProduct.Type = daoCoin.Type
	daoProduct.Address = daoCoin.TokenAddress
	daoProduct.ChainId = daoCoin.ChainId
	daoProduct.ChainName = daoCoin.ChainName
	daoProduct.Name = daoCoin.CoinName
	daoProduct.Symbol = strings.ToUpper(daoCoin.CoinSymbol)
	daoProduct.Category = `Crypto Projects` //Category for coin/ token
	daoProduct.Subcategory = daoCoin.Tag
	daoProduct.Image = daoCoin.CoinImage
	daoProduct.FromBy = `3rd` //from API 3rd
	daoProduct.Contract = daoCoin.Contract
	daoProduct.Description = daoCoin.Description
	daoProduct.Detail = daoCoin.Detail
}
