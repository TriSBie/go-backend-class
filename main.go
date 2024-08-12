package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"simple_bank.sqlc.dev/app/api"
	db "simple_bank.sqlc.dev/app/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	// get a database handle - using sql.Open to initialize the db variable
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}

	// initial store by using db.New
	store := db.NewStore(conn)
	server := api.NewServer(store)

	// cannot declare variable twice if had been using :=
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
