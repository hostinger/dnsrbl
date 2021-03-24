package hbl

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/labstack/echo/v4"
)

// @Summary     Block or Allow an IP address.
// @Description Use this endpoint to Block or Allow an IP address depending on Action argument in body.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200
// @Failure     422 {object} Error
// @Failure     500 {object} Error
// @Router      /addresses [POST]
func (api *API) handleAddressesPost(c echo.Context) error {
	var req BlockRequest
	var address Address
	if err := req.Bind(c, &address); err != nil {
		return echo.NewHTTPError(422, fmt.Sprintf("Failed to validate request body: %s", err))
	}
	a, err := api.Service.GetOne(context.Background(), address.IP)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	if a != nil {
		return echo.NewHTTPError(500, "Address already exists, delete it before other actions")
	}
	switch req.Action {
	case "Block":
		if err := api.Service.Block(context.Background(), &address); err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
		}
	case "Allow":
		if err := api.Service.Allow(context.Background(), &address); err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
		}
	}
	return c.JSON(200, nil)
}

// @Summary     Delete an IP address.
// @Description Use this endpoint to delete an already blocked or allowed IP address.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200
// @Param 		ip path string true "IP Address"
// @Failure     422 {object} Error
// @Failure     500 {object} Error
// @Router      /addresses/{ip} [DELETE]
func (api *API) handleAddressesDelete(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	address, err := api.Service.GetOne(context.Background(), ip)
	if err != nil && err == sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(404, "Address doesn't exist")
		}
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	if address.Action == "Block" {
		if err := api.Service.Unblock(context.Background(), address); err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
		}
	}
	if err := api.Service.Delete(context.Background(), ip); err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, nil)
}

// @Summary     Get an IP address.
// @Description Use this endpoint to fetch details about an already blocked or allowed IP address.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200 {object} Address
// @Param 		ip path string true "IP Address"
// @Failure     422 {object} Error
// @Failure     500 {object} Error
// @Router      /addresses/{ip} [GET]
func (api *API) handleAddressesGetOne(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	address, err := api.Service.GetOne(context.Background(), ip)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(404, "Address doesn't exist")
		}
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, address)
}

// @Summary     Get all IP addresses.
// @Description Use this endpoint to fetch details about all already blocked or allowed IP addresses.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200 {array} Address
// @Failure     422 {object} Error
// @Failure     500 {object} Error
// @Router      /addresses [GET]
func (api *API) handleAddressesGetAll(c echo.Context) error {
	addresses, err := api.Service.GetAll(context.Background())
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, addresses)
}

// @Summary     Get an IP address.
// @Description Use this endpoint to fetch details about an already blocked or allowed IP address.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200 {object} Address
// @Param 		name path string true "Name of the Checker"
// @Param 		ip path string true "IP Address"
// @Failure     422 {object} Error
// @Failure     500 {object} Error
// @Router      /addresses/check/{name}/{ip} [GET]
func (api *API) handleAddressesCheck(c echo.Context) error {
	name, ip := c.Param("name"), c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	result, err := api.Service.Check(context.Background(), name, ip)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, result)
}

func (api *API) handleAddressesSyncAll(c echo.Context) error {
	if err := api.Service.SyncAll(context.Background()); err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, nil)
}

func (api *API) handleAddressesSyncOne(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	if err := api.Service.SyncOne(context.Background(), ip); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(404, "Address doesn't exist")
		}
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, nil)
}

func (api *API) handleHealth(c echo.Context) error {
	return c.String(200, "OK")
}

func (api *API) handleVersion(c echo.Context) error {
	return c.String(200, "1.0.0")
}
