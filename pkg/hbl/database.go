package hbl

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL Drivers
)

func InitDB(username string, password string, host string, port string, database string, timeout time.Duration) (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, host, port, database)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("Failed to establish connection to the database after timeout: %s", timeout)
		case <-ticker.C:
			db, err := sql.Open("mysql", uri)
			if err != nil {
				return nil, err
			}
			err = db.Ping()
			if err == nil {
				return db, nil
			}
			log.Printf("Failed to establish connection to the database: %s", err)
		}
	}
}
