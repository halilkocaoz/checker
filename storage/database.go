package storage

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var connstr = os.Getenv("AZURE_POSTGRES_CONNSTR")

func UpsMoDBConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
