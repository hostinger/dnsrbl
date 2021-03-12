package hbl

import "github.com/labstack/echo/v4"

type Route struct {
	Method string
	Path   string
	Func   echo.HandlerFunc
}

func (api *API) GetRoutes() []*Route {
	return []*Route{
		// Common
		{
			Method: "GET",
			Path:   "/version",
			Func:   api.handleVersion,
		},
		{
			Method: "GET",
			Path:   "/health",
			Func:   api.handleHealth,
		},
		// Addresses
		{
			Method: "GET",
			Path:   "/api/v1/addresses",
			Func:   api.handleAddressesGetAll,
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses",
			Func:   api.handleAddressesPost,
		},
		{
			Method: "GET",
			Path:   "/api/v1/addresses/:ip",
			Func:   api.handleAddressesGetOne,
		},
		{
			Method: "DELETE",
			Path:   "/api/v1/addresses/:ip",
			Func:   api.handleAddressesDelete,
		},
		{
			Method: "GET",
			Path:   "/api/v1/addresses/check/:name/:ip",
			Func:   api.handleAddressesCheck,
		},
	}
}

func (api *API) SetupRoutes() {
	for _, route := range api.GetRoutes() {
		api.Server.Add(route.Method, route.Path, route.Func)
	}
}
