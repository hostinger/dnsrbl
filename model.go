package dnsrbl

import (
	"time"
)

// Address ...
type Address struct {
	Address   string
	Comment   string
	CreatedAt time.Time
}
