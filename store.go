package dnsrbl

import (
	"context"
	"database/sql"
)

type Store struct {
	Database *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Database: db,
	}
}

func (s *Store) CreateAddress(ctx context.Context, address Address) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "REPLACE INTO addresses(address, comment) VALUES (INET_ATON(?), ?)",
		address.Address, address.Comment)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteAddress(ctx context.Context, address string) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM addresses WHERE address = INET_ATON(?) LIMIT 1", address)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAddress(ctx context.Context, ip string) (address Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return Address{}, err
	}
	result := tx.QueryRowContext(ctx, "SELECT INET_NTOA(address), comment, created_at FROM addresses WHERE address = INET_ATON(?) LIMIT 1", ip)
	if err := result.Scan(&address.Address, &address.Comment, &address.CreatedAt); err != nil {
		return Address{}, err
	}

	if err := tx.Commit(); err != nil {
		return Address{}, err
	}
	return address, nil
}

func (s *Store) GetAddresses(ctx context.Context) (addresses []Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	results, err := tx.QueryContext(ctx, "SELECT INET_NTOA(address), comment, created_at FROM addresses")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for results.Next() {
		var address Address
		if err := results.Scan(&address.Address, &address.Comment, &address.CreatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return addresses, nil
}
