package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
