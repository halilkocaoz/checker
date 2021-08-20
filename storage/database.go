package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func UpsMoDBConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=UpsMo user=postgres password=password")

	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
