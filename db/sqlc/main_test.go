package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// declare global variables
var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error

	// get a database handle - using sql.Open to initialize the db variable
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
