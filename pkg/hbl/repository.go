package hbl

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetAddresses(
		ctx context.Context,
	) ([]*Address, error)

	CreateAddress(
		ctx context.Context, address *Address,
	) error

	DeleteAddress(
		ctx context.Context, ip string,
	) error

	GetAddress(
		ctx context.Context, ip string,
	) (*Address, error)
}

type MySQLRepository struct {
	DB *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{
		DB: db,
	}
}

func (s *MySQLRepository) CreateAddress(ctx context.Context, address *Address) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := `
		INSERT INTO
			addresses(
				ip,
				author,
				action,
				comment
			)
		VALUES
			(
				INET_ATON(?),
				?,
				?,
				?
			)
	`
	_, err = tx.ExecContext(ctx, q, address.IP, address.Author, address.Action, address.Comment)

	if err != nil {
		tx.Rollback() // nolint
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *MySQLRepository) DeleteAddress(ctx context.Context, ip string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := `
		DELETE FROM
			addresses
		WHERE
			ip = INET_ATON(?)
		LIMIT 1
	`
	_, err = tx.ExecContext(ctx, q, ip)
	if err != nil {
		tx.Rollback() // nolint
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *MySQLRepository) GetAddress(ctx context.Context, ip string) (*Address, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := `
		SELECT
			INET_NTOA(ip),
			author,
			action,
			comment,
			created_at
		FROM
			addresses
		WHERE
			ip = INET_ATON(?)
		LIMIT 1
	`
	var address Address
	result := tx.QueryRowContext(ctx, q, ip)
	if err := result.Scan(&address.IP, &address.Author, &address.Action,
		&address.Comment, &address.CreatedAt); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &address, nil
}

func (s *MySQLRepository) GetAddresses(ctx context.Context) ([]*Address, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := `
		SELECT
			INET_NTOA(ip),
			author,
			action,
			comment,
			created_at
		FROM
			addresses
	`
	results, err := tx.QueryContext(ctx, q)
	if err != nil {
		tx.Rollback() // nolint
		return nil, err
	}
	var addresses []*Address
	for results.Next() {
		var address Address
		if err := results.Scan(&address.IP, &address.Author, &address.Action,
			&address.Comment, &address.CreatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, &address)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return addresses, nil
}
