package hbl

import (
	"database/sql"
	"fmt"
	"log"
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
func (api *API) handleBlock(c echo.Context) error {
	var req BlockRequest
	var address Address
	if err := req.Bind(c, &address); err != nil {
		return c.JSON(422, ErrorResponse{Message: fmt.Sprintf("Failed to validate request body: %s", err)})
	}
	// Fetch metadata about the IP address
	report, err := api.abuseipdbClient.Check(req.IP)
	if err != nil {
		log.Printf("Failed to fetch AbuseIPDB metadata: %s", err)
	}
	address.Metadata.AbuseIpDbMetadata = AbuseIpDbMetadata{
		IP:                   req.IP,
		ISP:                  report.Data.Isp,
		UsageType:            report.Data.UsageType,
		CountryCode:          report.Data.CountryCode,
		TotalReports:         report.Data.TotalReports,
		LastReportedAt:       report.Data.LastReportedAt,
		NumDistinctUsers:     report.Data.NumDistinctUsers,
		AbuseConfidenceScore: report.Data.AbuseConfidenceScore,
	}
	// Create entity in database
	if err := api.Service.BlockAddress(address); err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute BlockAddress: %s", err)})
	}
	// Block in Cloudflare
	if *req.BlockCloudflare == true {
		if err := api.cfClient.BlockIPAddress(req.IP); err != nil {
			return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute BlockIPAddress: %s", err)})
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
func (api *API) handleListOne(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return c.JSON(422, ErrorResponse{Message: "Param 'ip' must be a valid IP address."})
	}
	address, err := api.Service.GetAddress(ip)
	if err != nil && err == sql.ErrNoRows {
		return c.JSON(404, ErrorResponse{Message: "Such IP address doesn't exist."})
	}
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute GetAddress: %s", err)})
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
func (api *API) handleUnblock(c echo.Context) error {
	ip := c.Param("ip")
	if net.ParseIP(ip) == nil {
		return c.JSON(422, ErrorResponse{Message: "Param 'ip' must be a valid IP address."})
	}
	if err := api.Service.UnblockAddress(ip); err != nil {
		return c.JSON(500, ErrorResponse{Message: fmt.Sprintf("Failed to execute UnblockAddress: %s", err)})
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
func (api *API) handleListAll(c echo.Context) error {
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
