package dnsrbl

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Address struct {
	Address   string
	Comment   string
	CreatedAt time.Time
}

func (m *Address) Validate() error {
	if net.ParseIP(m.Address) == nil {
		return errors.New("Field 'Address' must be a valid IP address.")
	}
	if len(strings.TrimSpace(m.Comment)) == 0 {
		return errors.New("Field 'Comment' must be not empty.")
	}
	return nil
}

type BlockAddressRequest struct {
	Address string `json:"address"`
	Comment string `json:"comment"`
}

func (m *BlockAddressRequest) Bind(c echo.Context, a *Address) error {
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := m.Validate(); err != nil {
		return err
	}
	a.Address = m.Address
	a.Comment = m.Comment
	return nil
}

func (m *BlockAddressRequest) Validate() error {
	if net.ParseIP(m.Address) == nil {
		return errors.New("Field 'Address' must be a valid IP address.")
	}
	if len(strings.TrimSpace(m.Comment)) == 0 {
		return errors.New("Field 'Comment' must be not empty.")
	}
	return nil
}
