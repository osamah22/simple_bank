package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	dbDriver := os.Getenv("DB_DRIVER")
	dbURL := os.Getenv("DB_URL")

	conn, err := sql.Open(dbDriver, dbURL)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}

	testDB = conn

	testQueries = New(conn)

	os.Exit(m.Run())
}
