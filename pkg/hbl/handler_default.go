package hbl

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/hostinger/hbl/pkg/logger"
	"github.com/labstack/echo/v4"
)

type handler struct {
	l       logger.Logger
	service Service
}

func NewDefaultHandler(l logger.Logger, s Service) Handler {
	return &handler{
		l:       l,
		service: s,
	}
}

// @Summary     Block or Allow an IP address.
// @Description Use this endpoint to Block or Allow an IP address depending on Action argument in body.
// @Produce     json
// @Accept      json
// @Tags        Addresses
// @Success     200
// @Router      /addresses [POST]
func (h *handler) HandleAddressesPost(c echo.Context) error {
	var req BlockRequest
	var address Address
	if err := req.Bind(c, &address); err != nil {
		return echo.NewHTTPError(422, fmt.Sprintf("Failed to validate request body: %s", err))
	}
	a, err := h.service.GetOne(context.Background(), address.IP)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	if a != nil {
		return echo.NewHTTPError(500, "Address already exists, delete it before other actions")
	}
	switch req.Action {
	case "Block":
		if err := h.service.Block(context.Background(), &address); err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
		}
	case "Allow":
		if err := h.service.Allow(context.Background(), &address); err != nil {
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
// @Router      /addresses/{ip} [DELETE]
func (h *handler) HandleAddressesDelete(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	address, err := h.service.GetOne(context.Background(), ip)
	if err != nil && err == sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(404, "Address doesn't exist")
		}
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	if address.Action == "Block" {
		if err := h.service.Unblock(context.Background(), address); err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
		}
	}
	if err := h.service.Delete(context.Background(), ip); err != nil {
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
// @Router      /addresses/{ip} [GET]
func (h *handler) HandleAddressesGetOne(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	address, err := h.service.GetOne(context.Background(), ip)
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
// @Router      /addresses [GET]
func (h *handler) HandleAddressesGetAll(c echo.Context) error {
	addresses, err := h.service.GetAll(context.Background())
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
// @Router      /addresses/check/{name}/{ip} [GET]
func (h *handler) HandleAddressesCheck(c echo.Context) error {
	name, ip := c.Param("name"), c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	result, err := h.service.Check(context.Background(), name, ip)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, result)
}

func (h *handler) HandleAddressesSyncAll(c echo.Context) error {
	if err := h.service.SyncAll(context.Background()); err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, nil)
}

func (h *handler) HandleAddressesSyncOne(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(422, "Param 'IP' must be a valid IP address")
	}
	if err := h.service.SyncOne(context.Background(), ip); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(404, "Address doesn't exist")
		}
		return echo.NewHTTPError(500, fmt.Sprintf("Error: %s", err))
	}
	return c.JSON(200, nil)
}

func (h *handler) HandleHealth(c echo.Context) error {
	return c.String(200, "OK")
}

func (h *handler) HandleVersion(c echo.Context) error {
	return c.String(200, "1.0.0")
}
