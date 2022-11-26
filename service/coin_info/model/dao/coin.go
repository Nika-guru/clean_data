package dao

import (
	"database/sql"
	"fmt"
	"reflect"
	"review-service/pkg/db"
	"review-service/pkg/log"
)

type CoinRepo struct {
	Coins []*Coin
}

type Coin struct {
	//Insert
	CoinID       *string //can be native coin or token coin
	CoinSymbol   *string
	CoinName     *string
	ChainName    *string //native coin
	TokenAddress any     //can null or primitive values(string)

	//Update
	CoinDescription *string
	CoinImage       *string
	TokenDecimal    *float64
}

type CoinDB struct {
	CoinID       sql.NullString
	CoinSymbol   sql.NullString
	CoinName     sql.NullString
	ChainName    sql.NullString
	TokenAddress sql.NullString

	CoinDescription sql.NullString
	CoinImage       sql.NullString
	TokenDecimal    sql.NullFloat64
}

func (coin *Coin) convertTo(coinDB *CoinDB) {
	//######### CoinID is not null value  #########
	if coin.CoinID != nil {
		coinDB.CoinID = sql.NullString{String: *coin.CoinID, Valid: true}
	} else {
		coinDB.CoinID = sql.NullString{Valid: false}
	}

	//######### CoinSymbol is not null value  #########
	if coin.CoinSymbol != nil {
		coinDB.CoinSymbol = sql.NullString{String: *coin.CoinSymbol, Valid: true}
	} else {
		coinDB.CoinSymbol = sql.NullString{Valid: false}
	}

	//######### CoinName is not null value  #########
	if coin.CoinName != nil {
		coinDB.CoinName = sql.NullString{String: *coin.CoinName, Valid: true}
	} else {
		coinDB.CoinName = sql.NullString{Valid: false}
	}

	//######### ChainName is not null value  #########
	if coin.ChainName != nil {
		coinDB.ChainName = sql.NullString{String: *coin.ChainName, Valid: true}
	} else {
		coinDB.ChainName = sql.NullString{Valid: false}
	}

	//######### TokenAddress is not null value  #########
	if coin.TokenAddress != nil {
		coinDB.TokenAddress = sql.NullString{String: coin.TokenAddress.(string), Valid: true}
	} else {
		coinDB.TokenAddress = sql.NullString{Valid: false}
	}

	//######### CoinDescription is not null value  #########
	if coin.CoinDescription != nil {
		coinDB.CoinDescription = sql.NullString{String: *coin.CoinDescription, Valid: true}
	} else {
		coinDB.CoinDescription = sql.NullString{Valid: false}
	}

	//######### CoinImage is not null value  #########
	if coin.CoinImage != nil {
		coinDB.CoinImage = sql.NullString{String: *coin.CoinImage, Valid: true}
	} else {
		coinDB.CoinImage = sql.NullString{Valid: false}
	}

	//######### TokenDecimal is not null value  #########
	if coin.TokenDecimal != nil {
		coinDB.TokenDecimal = sql.NullFloat64{Float64: *coin.TokenDecimal, Valid: true}
	} else {
		coinDB.TokenDecimal = sql.NullFloat64{Valid: false}
	}
}

func (coinDB *CoinDB) convertTo(coin *Coin) {
	if reflect.TypeOf(coinDB.CoinID) != nil {
		//######### CoinID is not null value#########
		if coinDB.CoinID.Valid {
			coin.CoinID = &coinDB.CoinID.String
		}
	}

	if reflect.TypeOf(coinDB.CoinSymbol) != nil {
		//######### CoinSymbol is not null value#########
		if coinDB.CoinSymbol.Valid {
			coin.CoinSymbol = &coinDB.CoinSymbol.String
		}
	}

	if reflect.TypeOf(coinDB.CoinName) != nil {
		//######### CoinName is not null value#########
		if coinDB.CoinName.Valid {
			coin.CoinName = &coinDB.CoinName.String
		}
	}

	if reflect.TypeOf(coinDB.ChainName) != nil {
		//######### ChainName is not null value#########
		if coinDB.ChainName.Valid {
			coin.ChainName = &coinDB.ChainName.String
		}
	}

	if reflect.TypeOf(coinDB.TokenAddress) != nil {
		//######### TokenAddress is not null value#########
		if coinDB.TokenAddress.Valid {
			coin.TokenAddress = coinDB.TokenAddress.String
		} else {
			coin.TokenAddress = nil
		}
	}

	if reflect.TypeOf(coinDB.CoinDescription) != nil {
		//######### CoinDescription is not null value#########
		if coinDB.CoinDescription.Valid {
			coin.CoinDescription = &coinDB.CoinDescription.String
		}
	}

	if reflect.TypeOf(coinDB.CoinImage) != nil {
		//######### CoinImage is not null value#########
		if coinDB.CoinImage.Valid {
			coin.CoinImage = &coinDB.CoinImage.String
		}
	}

	if reflect.TypeOf(coinDB.TokenDecimal) != nil {
		//######### TokenDecimal is not null value#########
		if coinDB.TokenDecimal.Valid {
			val := coinDB.TokenDecimal.Float64
			coin.TokenDecimal = &val
		}
	}
}

func (coin *Coin) GetKeyMap() string {
	coinKey := ""
	if coin.CoinID != nil {
		coinKey += *coin.CoinID
	}
	if coin.ChainName == nil {
		coinKey += fmt.Sprintf("%v%v", nil, coin.TokenAddress)
	} else {
		coinKey += fmt.Sprintf("%s%v", *coin.ChainName, coin.TokenAddress) //dont cast type(can't nil or string)
	}
	return coinKey
}

func (repo *CoinRepo) SelectAllCoinsTokens() error {
	query :=
		`
	SELECT coin_id, coin_symbol, coin_name,
			coin_description, coin_image, chain_name,
			token_address, token_decimals
	FROM coins
	--where coin_id = 'auxilium' --test
	--where coin_image is null -- coin chua cao duoc
	ORDER BY coin_id -- same coin_id group nearly
	`
	rows, err := db.PSQL.Query(query)
	if err != nil {
		log.Println(log.LogLevelError, "dao/coins.go/SelectAllCoins", err.Error())
		return err
	}

	defer rows.Close()

	for rows.Next() {
		coinDB := &CoinDB{}
		coinDAO := &Coin{}
		err = rows.Scan(
			&coinDB.CoinID, &coinDB.CoinSymbol, &coinDB.CoinName,
			&coinDB.CoinDescription, &coinDB.CoinImage, &coinDB.ChainName,
			&coinDB.TokenAddress, &coinDB.TokenDecimal)
		coinDB.convertTo(coinDAO)

		repo.Coins = append(repo.Coins, coinDAO)

		if err != nil {
			log.Println(log.LogLevelError, "dao/coins.go/SelectAllCoins", err.Error())
			return err
		}
	}

	return nil
}

func (repo *CoinRepo) InsertDB() (insertedRepo *CoinRepo) {
	existedCoinRepo := &CoinRepo{}
	//########## Get all items from database to exitedRepo ##########
	existedCoinRepo.SelectAllCoinsTokens()

	//########## Set init default ##########
	insertedRepo = &CoinRepo{}

	exitedCoinMap := make(map[string]bool)
	//########## Traverse each existed Coin in Database ##########
	for _, existedCoin := range existedCoinRepo.Coins {
		if existedCoin != nil {
			exitedCoinMap[existedCoin.GetKeyMap()] = true
		}
	}

	//########## Traverse each crawled Coin ##########
	for _, coin := range repo.Coins {
		_, found := exitedCoinMap[coin.GetKeyMap()]

		//########## Current crawled platform not existed in database ##########
		if !found {
			coinDB := &CoinDB{}
			coin.convertTo(coinDB)

			query := `
			INSERT INTO coins
					(coin_id, coin_symbol, coin_name, coin_description, coin_image, chain_name, token_address, token_decimals)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8);

			`
			_, err := db.PSQL.Exec(query,
				coinDB.CoinID, coinDB.CoinSymbol, coinDB.CoinName,
				coinDB.CoinDescription, coinDB.CoinImage, coinDB.ChainName,
				coinDB.TokenAddress, coinDB.TokenDecimal)

			//########## Insert database failed ##########
			if err != nil {
				log.Println(log.LogLevelError, "dao/coin.go/InsertDB", err.Error())
				continue
			} else
			//########## Insert database successful ##########
			{
				coinVal := coin
				insertedRepo.Coins = append(insertedRepo.Coins, coinVal)
			}

		}

	}

	return insertedRepo
}

func (repo *CoinRepo) UpdateDBDecimalImageDescription() (updatedTotal int) {
	existedCoinRepo := &CoinRepo{}
	//########## Get all items from database to exitedRepo ##########
	existedCoinRepo.SelectAllCoinsTokens()

	//########## Set Updated Total of query = 0 ##########
	updatedTotal = 0

	exitedCoinMap := make(map[string]Coin)
	//########## Traverse each existed Coin in Database ##########
	for _, existedCoin := range existedCoinRepo.Coins {
		if existedCoin != nil {
			exitedCoinMap[existedCoin.GetKeyMap()] = *existedCoin
		}
	}

	//########## Traverse each crawled Coin ##########
	for _, coin := range repo.Coins {
		_, found := exitedCoinMap[coin.GetKeyMap()]

		//########## Current crawled platform not existed in database ##########
		if found {
			coinDB := &CoinDB{}
			coin.convertTo(coinDB)

			query :=
				`
			UPDATE coins
				SET coin_description=$2, coin_image=$3, token_decimals=$4
				where coin_id = $1;
			`
			_, err := db.PSQL.Exec(query,
				coinDB.CoinID, coinDB.CoinDescription, coinDB.CoinImage,
				coinDB.TokenDecimal)

			//########## Insert database failed ##########
			if err != nil {
				log.Println(log.LogLevelError, "dao/coin.go/InsertDB", err.Error())
				continue
			} else
			//########## Insert database successful ##########
			{
				//test
				// fmt.Println("Updated ", coinDB.CoinID.String, coinDB.ChainName.String, coinDB.TokenAddress.String, coinDB.TokenDecimal.Float64)
				updatedTotal += 1
			}

		}

	}
	return updatedTotal
}
