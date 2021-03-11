package checkers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caarlos0/env"
)

type AbuseIPDBReport struct {
	IP                   string
	CountryCode          string
	UsageType            string
	ISP                  string
	AbuseConfidenceScore int
	NumDistinctUsers     int
	TotalReports         int
	LastReportedAt       *time.Time
}

type AbuseIPDBChecker struct {
	Client  *http.Client
	DB      *sql.DB
	BaseURL string
	Key     string `env:"ABUSEIPDB_API_KEY,required"`
}

func NewAbuseIPDBChecker(db *sql.DB) Checker {
	c := &AbuseIPDBChecker{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseURL: "https://api.abuseipdb.com/api/v2",
		DB:      db,
	}
	if err := env.Parse(c); err != nil {
		log.Fatalf("Checkers: AbuseIPDBChecker: Failed to initialize from environment variables: %s", err)
	}
	return c
}

func (c *AbuseIPDBChecker) Name() string {
	return "AbuseIPDB"
}

func (c *AbuseIPDBChecker) Call(ctx context.Context, ip string) (*AbuseIPDBReport, error) {
	type Result struct {
		Data struct {
			Hostnames            []interface{} `json:"hostnames"`
			IPAddress            string        `json:"ipAddress"`
			CountryCode          string        `json:"countryCode"`
			UsageType            string        `json:"usageType"`
			Isp                  string        `json:"isp"`
			Domain               string        `json:"domain"`
			IPVersion            int           `json:"ipVersion"`
			AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
			TotalReports         int           `json:"totalReports"`
			NumDistinctUsers     int           `json:"numDistinctUsers"`
			LastReportedAt       *time.Time    `json:"lastReportedAt"`
			IsWhitelisted        bool          `json:"isWhitelisted"`
			IsPublic             bool          `json:"isPublic"`
		} `json:"data"`
	}

	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("argument must be a valid IP address")
	}

	uri := fmt.Sprintf("%s/check", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("ipAddress", ip)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Key", c.Key)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failure from AbuseIPDB API: %s", string(body))
	}

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &AbuseIPDBReport{
		IP:                   result.Data.IPAddress,
		CountryCode:          result.Data.CountryCode,
		UsageType:            result.Data.UsageType,
		ISP:                  result.Data.Isp,
		AbuseConfidenceScore: result.Data.AbuseConfidenceScore,
		NumDistinctUsers:     result.Data.NumDistinctUsers,
		TotalReports:         result.Data.TotalReports,
		LastReportedAt:       result.Data.LastReportedAt,
	}, nil
}

func (c *AbuseIPDBChecker) Check(ctx context.Context, ip string) (interface{}, error) {
	report, err := c.GetReport(ctx, ip)
	if err != nil && err == sql.ErrNoRows {
		report, err = c.Call(ctx, ip)
		if err != nil {
			return nil, err
		}
		if err := c.SaveReport(ctx, report); err != nil {
			return nil, err
		}
	}
	return report, nil
}

func (c *AbuseIPDBChecker) SaveReport(ctx context.Context, report *AbuseIPDBReport) error {
	tx, err := c.DB.BeginTx(ctx, nil)
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
		report.IP, report.AbuseConfidenceScore,
		report.CountryCode, report.UsageType,
		report.ISP, report.TotalReports,
		report.NumDistinctUsers,
		report.LastReportedAt,
	)
	if err != nil {
		tx.Rollback() // nolint
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *AbuseIPDBChecker) GetReport(ctx context.Context, ip string) (*AbuseIPDBReport, error) {
	var report AbuseIPDBReport
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
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
	if err := result.Scan(&report.IP, &report.AbuseConfidenceScore,
		&report.CountryCode, &report.UsageType,
		&report.ISP, &report.TotalReports,
		&report.NumDistinctUsers,
		&report.LastReportedAt); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &report, nil
}
