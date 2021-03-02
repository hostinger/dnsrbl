package dnsrbl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/hostinger/dnsrbl/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type API struct {
	Config   *Config
	Server   *echo.Echo
	Database *sql.DB
	Service  *Service
}

func NewAPI(cfg *Config, db *sql.DB) *API {
	server := echo.New()

	server.HideBanner = true
	server.HidePort = true

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	store := NewStore(db)
	service := NewService(cfg, store)

	return &API{
		Service:  service,
		Server:   server,
		Config:   cfg,
		Database: db,
	}
}

func (api *API) init() {
	// Blocklist Routes
	{
		api.Server.Add("DELETE", "/api/v1/blocklist/:ip", api.handleBlocklistDelete)
		api.Server.Add("GET", "/api/v1/blocklist/:ip", api.handleBlocklistGet)
		api.Server.Add("GET", "/api/v1/blocklist", api.handleBlocklistGetAll)
		api.Server.Add("POST", "/api/v1/blocklist", api.handleBlocklistPost)
	}
	// Common Routes
	{
		api.Server.Add("GET", "/version", api.handleVersion)
		api.Server.Add("GET", "/health", api.handleHealth)
	}
	// Swagger
	{
		api.Server.Add("GET", "/swagger/*", echoSwagger.WrapHandler)
	}
}

func (api *API) Start(host string, port string) error {
	api.init()
	listenAddress := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on %s", listenAddress)
	if err := api.Server.Start(listenAddress); err != nil && err != http.ErrServerClosed {
		api.Server.Logger.Fatal("Exiting.")
	}
	return nil
}

func (api *API) Stop() {
	fmt.Println() // For CTRL+C
	log.Println("Exiting...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
	}()
	if err := api.Server.Shutdown(ctx); err != nil {
		log.Fatal("Failed clean shutdown, exiting with errors.")
	}
}
