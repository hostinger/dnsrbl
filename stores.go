package dnsrbl

import (
	"context"
	"database/sql"
)

// AddressStore ...
type AddressStore interface {
	Get(ctx context.Context, ip string) (address Address, err error)
	GetAll(ctx context.Context) (addresses []Address, err error)
	Create(ctx context.Context, address Address) (err error)
	Delete(ctx context.Context, address string) (err error)
}

// MySQLAddressStore ...
type MySQLAddressStore struct {
	Database *sql.DB
}

// NewMySQLAddressStore ...
func NewMySQLAddressStore(db *sql.DB) *MySQLAddressStore {
	return &MySQLAddressStore{
		Database: db,
	}
}

// Create ...
func (s MySQLAddressStore) Create(ctx context.Context, address Address) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO addresses(address, comment, expires_at) VALUES (INET_ATON(?), ?, ?)",
		address.Address, address.Comment, address.ExpiresAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Delete ...
func (s MySQLAddressStore) Delete(ctx context.Context, address string) error {
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

// Get ...
func (s MySQLAddressStore) Get(ctx context.Context, ip string) (address Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return Address{}, err
	}
	result := tx.QueryRowContext(ctx, "SELECT INET_NTOA(address), comment, expires_at FROM addresses WHERE address = INET_ATON(?) LIMIT 1", ip)
	if err := result.Scan(&address.Address, &address.Comment, &address.ExpiresAt); err != nil {
		return Address{}, err
	}

	if err := tx.Commit(); err != nil {
		return Address{}, err
	}
	return address, nil
}

// GetAll ...
func (s MySQLAddressStore) GetAll(ctx context.Context) (addresses []Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	results, err := tx.QueryContext(ctx, "SELECT INET_NTOA(address), comment, expires_at FROM addresses")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for results.Next() {
		var address Address
		if err := results.Scan(&address.Address, &address.Comment, &address.ExpiresAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return addresses, nil
}
