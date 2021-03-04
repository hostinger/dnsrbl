package hbl

import (
	"context"
	"database/sql"
)

type AddressStore struct {
	Database *sql.DB
}

func NewAddressStore(db *sql.DB) *AddressStore {
	return &AddressStore{
		Database: db,
	}
}

func (s *AddressStore) Create(ctx context.Context, address Address) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := `
		REPLACE INTO
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
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *AddressStore) Delete(ctx context.Context, ip string) error {
	tx, err := s.Database.BeginTx(ctx, nil)
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
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *AddressStore) GetOne(ctx context.Context, ip string) (address Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return Address{}, err
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
	result := tx.QueryRowContext(ctx, q, ip)
	if err := result.Scan(&address.IP, &address.Author, &address.Action, &address.Comment, &address.CreatedAt); err != nil {
		return Address{}, err
	}
	if err := tx.Commit(); err != nil {
		return Address{}, err
	}
	return address, nil
}

func (s *AddressStore) GetAll(ctx context.Context) (addresses []Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
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
		tx.Rollback()
		return nil, err
	}
	for results.Next() {
		var address Address
		if err := results.Scan(&address.IP, &address.Author, &address.Action, &address.Comment, &address.CreatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *AddressStore) GetAllByAction(ctx context.Context, action string) (addresses []Address, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
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
			action = ?
	`
	results, err := tx.QueryContext(ctx, q, action)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for results.Next() {
		var address Address
		if err := results.Scan(&address.IP, &address.Author, &address.Action, &address.Comment, &address.CreatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return addresses, nil
}

type MetadataStore struct {
	Database *sql.DB
}

func NewMetadataStore(db *sql.DB) *MetadataStore {
	return &MetadataStore{
		Database: db,
	}
}

func (s *MetadataStore) Create(ctx context.Context, metadata AbuseIpDbMetadata) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := `
		REPLACE INTO
			abuseipdb_metadata(
				ip,
				abuse_confidence_score,
				country_code,
				usage_type,
				isp,
				total_reports,
				num_distinct_users,
				last_reported_at
			)
		VALUES
			(
				INET_ATON(?),
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)
	`
	_, err = tx.ExecContext(ctx, q,
		&metadata.IP, &metadata.AbuseConfidenceScore,
		&metadata.CountryCode, &metadata.UsageType,
		&metadata.ISP, &metadata.TotalReports,
		&metadata.NumDistinctUsers,
		&metadata.LastReportedAt,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *MetadataStore) Delete(ctx context.Context, ip string) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := `
		DELETE FROM
			abuseipdb_metadata
		WHERE
			ip = INET_ATON(?)
		LIMIT 1
	`
	_, err = tx.ExecContext(ctx, q, ip)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *MetadataStore) GetOne(ctx context.Context, ip string) (metadata AbuseIpDbMetadata, err error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return AbuseIpDbMetadata{}, err
	}
	q := `
		SELECT
			INET_NTOA(ip),
			abuse_confidence_score,
			country_code,
			usage_type,
			isp,
			total_reports,
			num_distinct_users,
			last_reported_at
		FROM
			abuseipdb_metadata
		WHERE
			ip = INET_ATON(?)
		LIMIT 1
	`
	result := tx.QueryRowContext(ctx, q, ip)
	if err := result.Scan(&metadata.IP, &metadata.AbuseConfidenceScore,
		&metadata.CountryCode, &metadata.UsageType,
		&metadata.ISP, &metadata.TotalReports,
		&metadata.NumDistinctUsers,
		&metadata.LastReportedAt); err != nil {
		return AbuseIpDbMetadata{}, err
	}
	if err := tx.Commit(); err != nil {
		return AbuseIpDbMetadata{}, err
	}
	return metadata, nil
}
