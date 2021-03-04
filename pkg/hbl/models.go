package hbl

import (
	"time"
)

type ErrorResponse struct {
	Message string
}

type AbuseIpDbMetadata struct {
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
	IP                  string
	Author              string
	Comment             string
	IsBlockedPDNS       bool
	IsBlockedCloudflare bool
	CreatedAt           time.Time
	Metadata            AddressMetadata
}

type AddressMetadata struct {
	AbuseIpDbMetadata AbuseIpDbMetadata
}
