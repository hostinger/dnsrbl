package hbl

import (
	"time"
)

type Address struct {
	IP        string
	Author    string
	Action    string
	Comment   string
	CreatedAt time.Time
}
