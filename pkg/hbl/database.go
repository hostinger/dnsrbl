package hbl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL Drivers
	"github.com/pkg/errors"
)

func InitDB(ctx context.Context, username, password, host, port, database string) (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, host, port, database)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout exceeded")
		case <-ticker.C:
			db, err := sql.Open("mysql", uri)
			if err != nil {
				return nil, err
			}
			err = db.Ping()
			if err == nil {
				return db, nil
			}
			log.Printf("failed to establish connection to the database: %s", err)
		}
	}
}
