package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL Drivers
)

func Init(ctx context.Context, username, password, host, port, db string) (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		db,
	)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("Timeout exceeded")
		case <-ticker.C:
			db, err := sql.Open("mysql", uri)
			if err != nil {
				return nil, err
			}
			err = db.Ping()
			if err == nil {
				return db, nil
			}
			log.Printf("Database: %s", err)
		}
	}
}
