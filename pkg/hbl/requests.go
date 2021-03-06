package hbl

import (
	"errors"
	"net"
	"strings"

	"github.com/labstack/echo/v4"
)

type BlockRequest struct {
	IP      string
	Author  string
	Action  string
	Comment string
}

func (m *BlockRequest) Bind(c echo.Context, a *Address) error {
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := m.Validate(); err != nil {
		return err
	}
	a.IP = m.IP
	a.Action = m.Action
	a.Author = m.Author
	a.Comment = m.Comment
	return nil
}

func (m *BlockRequest) Validate() error {
	if net.ParseIP(m.IP) == nil {
		return errors.New("Field 'IP' must be a valid IP address")
	}
	if strings.TrimSpace(m.Author) == "" {
		return errors.New("Field 'Author' must not be empty")
	}
	if strings.TrimSpace(m.Comment) == "" {
		return errors.New("Field 'Comment' must not be empty")
	}
	if strings.TrimSpace(m.Action) == "" {
		return errors.New("Field 'Action' must not be empty")
	}
	if m.Action != "Block" && m.Action != "Allow" {
		return errors.New("Field 'Action' must be valid")
	}
	return nil
}
