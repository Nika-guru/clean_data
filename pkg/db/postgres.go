package db

import (
	"database/sql"
	"fmt"

	"base/pkg/log"

	_ "github.com/lib/pq"
)

// PSQL Configuration Struct
type psqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// PSQL Configuration Variable
var psqlCfg psqlConfig

// MySQL Variable
var PSQL *sql.DB

// PSQL Connect Function
func psqlConnect() *sql.DB {
	// Initialize Connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", psqlCfg.Host, psqlCfg.Port, psqlCfg.User, psqlCfg.Password, psqlCfg.Name)
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(log.LogLevelFatal, "psql-connect", err.Error())
	}

	// Test Connection
	err = conn.Ping()
	if err != nil {
		log.Println(log.LogLevelFatal, "psql-connect", err.Error())
	}

	// Return Connection
	return conn
}
