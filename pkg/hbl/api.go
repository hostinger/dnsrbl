package hbl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/hostinger/hbl/docs" // Needed for Swagger
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type Config struct {
	DB     *sql.DB
	Logger *zap.Logger
	Host   string
	Port   string
}

type API struct {
	Cfg     *Config
	Server  *echo.Echo
	Service Service
}

func NewAPI(cfg *Config) *API {
	server := echo.New()
	server.HidePort = true
	server.HideBanner = true

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/swagger/*", echoSwagger.WrapHandler)

	service := NewService(cfg.Logger, NewMySQLRepository(cfg.Logger, cfg.DB))

	return &API{
		Cfg:     cfg,
		Server:  server,
		Service: service,
	}
}

func (api *API) Start() error {
	api.Cfg.Logger.Info("Starting", zap.String("host", api.Cfg.Host), zap.String("port", api.Cfg.Port))
	api.SetupRoutes()
	uri := fmt.Sprintf("%s:%s", api.Cfg.Host, api.Cfg.Port)
	if err := api.Server.Start(uri); err != nil && err != http.ErrServerClosed {
		api.Cfg.Logger.Fatal(err.Error())
	}
	return nil
}

func (api *API) Stop() {
	api.Cfg.Logger.Info("Stopping", zap.String("host", api.Cfg.Host), zap.String("port", api.Cfg.Port))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
	}()
	if err := api.Server.Shutdown(ctx); err != nil {
		log.Printf("Failed clean shutdown, exiting with errors.")
	}
}
