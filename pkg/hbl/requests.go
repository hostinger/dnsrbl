package hbl

import (
	"errors"
	"net"
	"strings"

	"github.com/labstack/echo/v4"
)

type BlockRequest struct {
	IP              string
	Author          string
	Comment         string
	BlockPDNS       *bool
	BlockCloudflare *bool
}

func (m *BlockRequest) Bind(c echo.Context, a *Address) error {
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := m.Validate(); err != nil {
		return err
	}
	a.IP = m.IP
	a.Author = m.Author
	a.Comment = m.Comment
	a.IsBlockedCloudflare = *m.BlockCloudflare
	a.IsBlockedPDNS = *m.BlockPDNS
	return nil
}

func (m *BlockRequest) Validate() error {
	if net.ParseIP(m.IP) == nil {
		return errors.New("Field 'IP' must be a valid IP address.")
	}
	if len(strings.TrimSpace(m.Comment)) == 0 {
		return errors.New("Field 'Comment' must not be empty.")
	}
	if m.BlockCloudflare == nil {
		return errors.New("Field 'BlockCloudflare' must not be empty.")
	}
	if m.BlockPDNS == nil {
		return errors.New("Field 'BlockPDNS' must not be empty.")
	}
	return nil
}
