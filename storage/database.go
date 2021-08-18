package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func UpMoDBConnection() (*sql.DB, error) {
	return sql.Open("postgres", "host=localhost port=5432 dbname=UpMo user=postgres password=password")
}
