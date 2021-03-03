package hbl

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL Drivers
)

func InitDB(username string, password string, host string, port string, database string) (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, host, port, database)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
