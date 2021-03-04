package hbl

import (
	"fmt"
	"net"

	"github.com/labstack/echo/v4"
)

// @Summary     Block an IP address
// @Description Block an IP address.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200
// @Failure     422 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /blocklist [POST]
func (api *API) handleAddressesPost(c echo.Context) error {
	var req BlockRequest
	var address Address
	if err := req.Bind(c, &address); err != nil {
		return c.JSON(422, ErrorResponse{Message: fmt.Sprintf("Failed to validate request body: %s", err)})
	}
	switch req.Action {
	case ActionBlock:
		if err := api.Service.BlockAddress(&address); err != nil {
			return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute Service.BlockAddress: %s", err)})
		}
	case ActionAllow:
		if err := api.Service.AllowAddress(&address); err != nil {
			return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute Service.AllowAddress: %s", err)})
		}
	}
	return c.JSON(200, nil)
}

// @Summary     Get a blocked IP address.
// @Description Get a blocked IP address.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Param 		ip path string true "Valid IP address to search for."
// @Success     200
// @Failure     422 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /blocklist/{ip} [GET]
func (api *API) handleAddressesGetOne(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return c.JSON(422, ErrorResponse{Message: "Param 'ip' must be a valid IP address."})
	}
	address, err := api.Service.GetAddress(ip)
	if err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute Service.GetAddress: %s", err)})
	}
	return c.JSON(200, address)
}

// @Summary     Delete a blocked IP address.
// @Description Delete a blocked IP address.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Param 		ip path string true "Valid IP address to search for."
// @Success     200
// @Failure     422 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /blocklist/{ip} [DELETE]
func (api *API) handleAddressesDelete(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return c.JSON(422, ErrorResponse{Message: "Param 'ip' must be a valid IP address."})
	}
	if err := api.Service.DeleteAddress(ip); err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute Service.DeleteAddress: %s", err)})
	}
	return c.JSON(200, nil)
}

// @Summary     Get all blocked IP addresses.
// @Description Get all blocked IP addresses.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200 {array} Address
// @Failure     422 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /blocklist/{ip} [DELETE]
func (api *API) handleAddressesGetAll(c echo.Context) error {
	addresses, err := api.Service.GetAddresses()
	if err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute GetAddresses: %s", err)})
	}
	return c.JSON(200, addresses)
}

func (api *API) handleHealth(c echo.Context) error {
	return c.String(200, "OK")
}

func (api *API) handleVersion(c echo.Context) error {
	return c.String(200, "1.0.0")
}
