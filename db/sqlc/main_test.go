package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"github.com/12138mICHAEL1111/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB


func TestMain(m *testing.M) { //will execute first
	config,err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource) //connect to database
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB) //testQueries is the db, but with less interface methods
	os.Exit(m.Run()) // m.RUN runs tests in account_test
}
