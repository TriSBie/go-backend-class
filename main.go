package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"simple_bank.sqlc.dev/app/api"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	// get a database handle - using sql.Open to initialize the db variable
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}

	// initial store by using db.New
	store := db.NewStore(conn)
	server := api.NewServer(store)

	// cannot declare variable twice if had been using :=
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
