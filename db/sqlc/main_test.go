package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/simpleBank/util"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("con not load configuration", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("connection failed")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}