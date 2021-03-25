package hbl

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Repository interface {
	GetAddress(ctx context.Context, ip string) (*Address, error)
	CreateAddress(ctx context.Context, address *Address) error
	GetAddresses(ctx context.Context) ([]*Address, error)
	DeleteAddress(ctx context.Context, ip string) error
}

type mysqlRepository struct {
	l  *zap.Logger
	DB *sql.DB
}

func NewMySQLRepository(l *zap.Logger, db *sql.DB) Repository {
	return &mysqlRepository{
		l:  l,
		DB: db,
	}
}

func (s *mysqlRepository) CreateAddress(ctx context.Context, address *Address) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.l.Error(
			"Failed to execute BeginTx",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "GetAddresses"),
			zap.Error(err),
		)
		return errors.Wrap(err, "Failed to execute BeginTx")
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
		s.l.Error(
			"Failed to execute ExecContext",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "CreateAddress"),
			zap.Error(err),
		)
		tx.Rollback() // nolint
		return errors.Wrap(err, "Failed to execute ExecContext")
	}

	if err := tx.Commit(); err != nil {
		s.l.Error(
			"Failed to execute Commit",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "CreateAddress"),
			zap.Error(err),
		)
		return errors.Wrap(err, "Failed to execute Commit")
	}
	return nil
}

func (s *mysqlRepository) DeleteAddress(ctx context.Context, ip string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.l.Error(
			"Failed to execute BeginTx",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "DeleteAddress"),
			zap.Error(err),
		)
		return errors.Wrap(err, "Failed to execute BeginTx")
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
		s.l.Error(
			"Failed to execute ExecContext",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "DeleteAddress"),
			zap.Error(err),
		)
		tx.Rollback() // nolint
		return errors.Wrap(err, "Failed to execute ExecContext")
	}
	if err := tx.Commit(); err != nil {
		s.l.Error(
			"Failed to execute Commit",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "DeleteAddress"),
			zap.Error(err),
		)
		return errors.Wrap(err, "Failed to execute Commit")
	}
	return nil
}

func (s *mysqlRepository) GetAddress(ctx context.Context, ip string) (*Address, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.l.Error(
			"Failed to execute BeginTx",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "GetAddress"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed to execute BeginTx")
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
		s.l.Error(
			"Failed to execute Commit",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "GetAddresses"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed to execute Commit")
	}
	return &address, nil
}

func (s *mysqlRepository) GetAddresses(ctx context.Context) ([]*Address, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.l.Error(
			"Failed to execute BeginTx",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "GetAddresses"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed to execute BeginTx")
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
		s.l.Error(
			"Failed to execute QueryContext",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "GetAddresses"),
			zap.Error(err),
		)
		tx.Rollback() // nolint
		return nil, errors.Wrap(err, "Failed to execute QueryContext")
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
		s.l.Error(
			"Failed to execute Commit",
			zap.String("repository", "MySQLRepository"),
			zap.String("method", "GetAddresses"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed to execute Commit")
	}
	return addresses, nil
}
