package hbl

import (
	"time"
)

const (
	ActionBlock = "Block"
	ActionAllow = "Allow"
)

type ErrorResponse struct {
	Message string
}

type AbuseIPDBMetadata struct {
	IP                   string
	ISP                  string
	UsageType            string
	CountryCode          string
	TotalReports         int
	NumDistinctUsers     int
	AbuseConfidenceScore int
	LastReportedAt       *time.Time
}

type Address struct {
	IP        string
	Author    string
	Action    string
	Comment   string
	CreatedAt time.Time
	Metadata  AddressMetadata
}

type AddressMetadata struct {
	AbuseIPDBMetadata AbuseIPDBMetadata
}
