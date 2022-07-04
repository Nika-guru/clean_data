package main

import (
	"log"
	"strings"

	"base/pkg/db"
	"base/pkg/router"
	"base/pkg/server"

	"base/service"
)

// Server Variable
var svr *server.Server

// Init Function
func init() {
	// Set Go Log Flags
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Load Routes
	service.LoadRoutes()

	// Initialize Server
	svr = server.NewServer(router.Router)
}

// Main Function
func main() {
	// Starting Server
	svr.Start()

	// Stopping Server
	defer svr.Stop()

	// Close Any Database Connections
	if len(server.Config.GetString("DB_DRIVER")) != 0 {
		switch strings.ToLower(server.Config.GetString("DB_DRIVER")) {
		case "postgres":
			defer db.PSQL.Close()
		case "mysql":
			defer db.MySQL.Close()
		case "mongo":
			defer db.MongoSession.Close()
		}
	}

}
