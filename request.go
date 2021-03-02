package dnsrbl

import "github.com/labstack/echo/v4"

type blockAddressRequest struct {
	Address string `json:"address" validate:"required"`
	Comment string `json:"comment" validate:"required"`
}

func (r *blockAddressRequest) bind(c echo.Context, m *Address) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	m.Address = r.Address
	m.Comment = r.Comment
	return nil
}
