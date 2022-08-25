package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) { //will execute first
	db, err := sql.Open(dbDriver, dbSource) //connect to database
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(db) //testQueries is the db, but with less interface methods
	os.Exit(m.Run()) // m.RUN runs tests in account_test
}
