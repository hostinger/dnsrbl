package hbl

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

	q := "REPLACE INTO addresses(ip, author, comment, is_blocked_pdns, is_blocked_cloudflare) VALUES (INET_ATON(?), ?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, q, address.IP, address.Author, address.Comment, address.IsBlockedPDNS, address.IsBlockedCloudflare)

	q = "REPLACE INTO abuseipdb_metadata(ip, abuse_confidence_score, country_code, usage_type, isp, total_reports, num_distinct_users, last_reported_at) VALUES (INET_ATON(?), ?, ?, ?, ?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, q,
		&address.Metadata.AbuseIpDbMetadata.IP, &address.Metadata.AbuseIpDbMetadata.AbuseConfidenceScore,
		&address.Metadata.AbuseIpDbMetadata.CountryCode, &address.Metadata.AbuseIpDbMetadata.UsageType,
		&address.Metadata.AbuseIpDbMetadata.ISP, &address.Metadata.AbuseIpDbMetadata.TotalReports,
		&address.Metadata.AbuseIpDbMetadata.NumDistinctUsers,
		&address.Metadata.AbuseIpDbMetadata.LastReportedAt,
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

func (s *Store) DeleteAddress(ctx context.Context, ip string) error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM addresses WHERE ip = INET_ATON(?) LIMIT 1", ip)
	_, err = tx.ExecContext(ctx, "DELETE FROM abuseipdb_metadata WHERE ip = INET_ATON(?) LIMIT 1", ip)
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
	q := `
		SELECT
			INET_NTOA(addresses.ip),
			addresses.author,
			addresses.comment,
			addresses.is_blocked_pdns,
			addresses.is_blocked_cloudflare,
			addresses.created_at,
			INET_NTOA(abuseipdb.ip),
			abuseipdb.abuse_confidence_score,
			abuseipdb.country_code,
			abuseipdb.usage_type,
			abuseipdb.isp,
			abuseipdb.total_reports,
			abuseipdb.num_distinct_users,
			abuseipdb.last_reported_at
		FROM
			addresses
		AS
			addresses
		INNER JOIN
			abuseipdb_metadata
		AS
			abuseipdb
		ON
			addresses.ip=abuseipdb.ip
		WHERE
			addresses.ip = INET_ATON(?)
		LIMIT 1
	`
	result := tx.QueryRowContext(ctx, q, ip)
	if err := result.Scan(&address.IP, &address.Author, &address.Comment, &address.IsBlockedPDNS, &address.IsBlockedCloudflare, &address.CreatedAt,
		&address.Metadata.AbuseIpDbMetadata.IP, &address.Metadata.AbuseIpDbMetadata.AbuseConfidenceScore,
		&address.Metadata.AbuseIpDbMetadata.CountryCode, &address.Metadata.AbuseIpDbMetadata.UsageType,
		&address.Metadata.AbuseIpDbMetadata.ISP, &address.Metadata.AbuseIpDbMetadata.TotalReports,
		&address.Metadata.AbuseIpDbMetadata.NumDistinctUsers,
		&address.Metadata.AbuseIpDbMetadata.LastReportedAt,
	); err != nil {
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
	q := `
		SELECT
			INET_NTOA(addresses.ip),
			addresses.author,
			addresses.comment,
			addresses.is_blocked_pdns,
			addresses.is_blocked_cloudflare,
			addresses.created_at,
			INET_NTOA(abuseipdb.ip),
			abuseipdb.abuse_confidence_score,
			abuseipdb.country_code,
			abuseipdb.usage_type,
			abuseipdb.isp,
			abuseipdb.total_reports,
			abuseipdb.num_distinct_users,
			abuseipdb.last_reported_at
		FROM
			addresses
		AS
			addresses
		INNER JOIN
			abuseipdb_metadata
		AS
			abuseipdb
		ON
			addresses.ip=abuseipdb.ip
	`
	results, err := tx.QueryContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for results.Next() {
		var address Address
		if err := results.Scan(&address.IP, &address.Author, &address.Comment, &address.IsBlockedPDNS, &address.IsBlockedCloudflare, &address.CreatedAt,
			&address.Metadata.AbuseIpDbMetadata.IP, &address.Metadata.AbuseIpDbMetadata.AbuseConfidenceScore,
			&address.Metadata.AbuseIpDbMetadata.CountryCode, &address.Metadata.AbuseIpDbMetadata.UsageType,
			&address.Metadata.AbuseIpDbMetadata.ISP, &address.Metadata.AbuseIpDbMetadata.TotalReports,
			&address.Metadata.AbuseIpDbMetadata.NumDistinctUsers,
			&address.Metadata.AbuseIpDbMetadata.LastReportedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return addresses, nil
}
