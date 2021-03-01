package dnsrbl

import (
	"time"
)

// Address ...
type Address struct {
	Address   string
	Comment   string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
