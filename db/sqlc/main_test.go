package db

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

func init() {
	// Load .env file if it exists, but ignore error if file is missing
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found or failed to load, continuing...")
	}
}

func TestMain(m *testing.M) {
	var dbDriver, dbURL string
	if err := godotenv.Load("../.env"); err != nil {
		dbDriver = "postgres"
		dbURL = "postgres://env_user:secret@localhost:5432/simple_bank?sslmode=disable"
	} else {
		dbDriver = os.Getenv("DB_DRIVER")
		dbURL = os.Getenv("DB_URL")
	}

	conn, err := sql.Open(dbDriver, dbURL)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}

	testDB = conn

	testQueries = New(conn)

	os.Exit(m.Run())
}
