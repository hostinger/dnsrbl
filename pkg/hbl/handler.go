package hbl

import "github.com/labstack/echo/v4"

type Handler interface {
	HandleHealth(c echo.Context) error
	HandleVersion(c echo.Context) error
	HandleAddressesPost(c echo.Context) error
	HandleAddressesCheck(c echo.Context) error
	HandleAddressesGetOne(c echo.Context) error
	HandleAddressesGetAll(c echo.Context) error
	HandleAddressesDelete(c echo.Context) error
	HandleAddressesSyncOne(c echo.Context) error
	HandleAddressesSyncAll(c echo.Context) error
}
