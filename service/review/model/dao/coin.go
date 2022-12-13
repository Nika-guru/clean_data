package dao

type Coin struct {
	CoinId    string
	Type      string //Token ^ Coin
	Address   string //NULL: 'NULL', lower
	ChainId   string //NULL: 'NULL'
	ChainName string //NULL: 'NULL'
	Symbol    string //Upper
	Name      string
	Tag       string

	//Dynamic
	TotalSupply    string
	MaxSupply      string
	Marketcap      string
	VolumneTrading string

	Image   string
	Decimal int //NULL: 0
	Src     string
	Detail  map[string]any
}

func (dao *Coin) InsertDB() error {
	return nil
}
