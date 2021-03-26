package hbl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/hostinger/hbl/docs" // Needed for Swagger
	"github.com/hostinger/hbl/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type Config struct {
	Handler Handler
	Logger  logger.Logger
	Host    string
	Port    string
}

type API struct {
	Handler Handler
	Cfg     *Config
	Server  *echo.Echo
}

func NewAPI(cfg *Config) *API {
	server := echo.New()
	server.HidePort = true
	server.HideBanner = true

	server.Use(middleware.Recover())

	server.GET("/swagger/*", echoSwagger.WrapHandler)

	return &API{
		Cfg:     cfg,
		Handler: cfg.Handler,
		Server:  server,
	}
}

func (api *API) Initialize(ctx context.Context) {
	api.SetupRoutes()
}

func (api *API) Start() {
	api.Cfg.Logger.Info("Starting", zap.String("host", api.Cfg.Host), zap.String("port", api.Cfg.Port))
	uri := fmt.Sprintf("%s:%s", api.Cfg.Host, api.Cfg.Port)
	if err := api.Server.Start(uri); err != nil && err != http.ErrServerClosed {
		api.Cfg.Logger.Fatal(
			"Failed to start, exiting with errors",
			zap.String("host", api.Cfg.Host),
			zap.String("port", api.Cfg.Port),
			zap.Error(err),
		)
	}
}

func (api *API) Stop() {
	api.Cfg.Logger.Info("Stopping", zap.String("host", api.Cfg.Host), zap.String("port", api.Cfg.Port))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
	}()
	if err := api.Server.Shutdown(ctx); err != nil {
		api.Cfg.Logger.Fatal(
			"Failed to stop, exiting with errors",
			zap.String("host", api.Cfg.Host),
			zap.String("port", api.Cfg.Port),
			zap.Error(err),
		)
	}
}
