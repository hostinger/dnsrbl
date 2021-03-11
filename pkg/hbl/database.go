package hbl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL Drivers
	"github.com/pkg/errors"
)

func InitDB(ctx context.Context) (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB"),
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
