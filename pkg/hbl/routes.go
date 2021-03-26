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
			Func:   api.Handler.HandleVersion,
		},
		{
			Method: "GET",
			Path:   "/health",
			Func:   api.Handler.HandleHealth,
		},
		// Addresses
		{
			Method: "GET",
			Path:   "/api/v1/addresses",
			Func:   api.Handler.HandleAddressesGetAll,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses",
			Func:   api.Handler.HandleAddressesPost,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "GET",
			Path:   "/api/v1/addresses/:ip",
			Func:   api.Handler.HandleAddressesGetOne,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "DELETE",
			Path:   "/api/v1/addresses/:ip",
			Func:   api.Handler.HandleAddressesDelete,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "GET",
			Path:   "/api/v1/addresses/check/:name/:ip",
			Func:   api.Handler.HandleAddressesCheck,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses/sync",
			Func:   api.Handler.HandleAddressesSyncAll,
			Middleware: []echo.MiddlewareFunc{
				KeyAuthMiddleware,
			},
		},
		{
			Method: "POST",
			Path:   "/api/v1/addresses/sync/:ip",
			Func:   api.Handler.HandleAddressesSyncOne,
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
