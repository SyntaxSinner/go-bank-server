package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/SyntaxSinner/BankCRUD_API/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var test_queries *sqlc.Queries
var test_db *sql.DB

func TestMain(m *testing.M) {
	var err error
	test_db, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Could not coonect to the database:", err)
	}

	test_queries = sqlc.New(test_db)

	os.Exit(m.Run())
}
