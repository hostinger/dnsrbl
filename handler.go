package dnsrbl

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/hostinger/dnsrbl/database"
	"github.com/labstack/echo/v4"
)

// @Summary Block an IP addresss
// @Description Block an IP address
// @Accept  json
// @Produce  json
// @Tags Addresses
// @Success 200
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /blocklist [POST]
func BlockHandler(c echo.Context) error {
	var a Address
	var req blockAddressRequest
	if err := req.bind(c, &a); err != nil {
		return c.JSON(422, ErrorResponse{Message: fmt.Sprintf("Failed to validate request body: %s", err)})
	}
	if err := BlockAddress(context.Background(), database.DB, a); err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute BlockAddress: %s", err)})
	}
	return c.JSON(200, nil)
}

// @Summary Unblock an IP addresss
// @Description Unblock an IP address
// @Accept  json
// @Produce  json
// @Param address path string true "The IP address"
// @Tags Addresses
// @Success 200
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /blocklist/{address} [DELETE]
func UnblockHandler(c echo.Context) error {
	address := c.Param("address")
	if net.ParseIP(address) == nil {
		return c.JSON(422, ErrorResponse{Message: "Param 'Address' must be a valid IP address."})
	}
	if err := UnblockAddress(context.Background(), database.DB, address); err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute UnblockAddress: %s", err)})
	}
	return c.JSON(200, nil)
}

// @Summary Get an IP address
// @Description Get an IP address
// @Accept  json
// @Produce  json
// @Param address path string true "The IP address"
// @Tags Addresses
// @Success 200 {object} Address
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /blocklist/{address} [GET]
func GetHandler(c echo.Context) error {
	ip := c.Param("address")
	if net.ParseIP(ip) == nil {
		return c.JSON(422, ErrorResponse{Message: "Param 'Address' must be a valid IP address."})
	}
	address, err := GetAddress(context.Background(), database.DB, ip)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute GetAddress: %s", err)})
	}
	if err != nil && err == sql.ErrNoRows {
		return c.JSON(404, nil)
	}
	return c.JSON(200, address)
}

// @Summary Get all IP addresses
// @Description Get all IP addresses
// @Accept  json
// @Produce  json
// @Tags Addresses
// @Param address path string true "The IP address"
// @Success 200 {array} Address
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /blocklist [GET]
func GetAllHandler(c echo.Context) error {
	addresses, err := GetAddresses(context.Background(), database.DB)
	if err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute GetAddresses: %s", err)})
	}
	return c.JSON(200, addresses)
}
