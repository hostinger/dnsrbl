package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL Drivers
)

// DB represents the database object.
var DB *sql.DB

// Init initializes the database or returns an error on failure.
func Init(username string, password string, host string, port string, database string) (err error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, host, port, database)
	DB, err = sql.Open("mysql", uri)
	if err != nil {
		return err
	}
	if err := DB.Ping(); err != nil {
		return err
	}
	return nil
}
