package hbl

import "github.com/labstack/echo/v4"

type Route struct {
	Method     string
	Path       string
	Func       echo.HandlerFunc
	Middleware []echo.MiddlewareFunc
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
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses",
			Func:   api.handleAddressesPost,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "GET",
			Path:   "/api/v1/addresses/:ip",
			Func:   api.handleAddressesGetOne,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "DELETE",
			Path:   "/api/v1/addresses/:ip",
			Func:   api.handleAddressesDelete,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "GET",
			Path:   "/api/v1/addresses/check/:name/:ip",
			Func:   api.handleAddressesCheck,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses/sync",
			Func:   api.handleAddressesSyncAll,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses/sync/:ip",
			Func:   api.handleAddressesSyncOne,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
	}
}

func (api *API) SetupRoutes() {
	for _, route := range api.GetRoutes() {
		api.Server.Add(route.Method, route.Path, route.Func, route.Middleware...)
	}
}
