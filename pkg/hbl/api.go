package hbl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/hostinger/dnsrbl/docs"
	"github.com/hostinger/dnsrbl/pkg/abuseipdb"
	"github.com/hostinger/dnsrbl/pkg/cloudflare"
	"github.com/hostinger/dnsrbl/pkg/dns"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type API struct {
	Server          *echo.Echo
	Database        *sql.DB
	Service         *Service
	cfClient        *cloudflare.Client
	abuseipdbClient *abuseipdb.Client
	dnsClient       *dns.Client
}

func NewAPI(db *sql.DB,
	cfClient *cloudflare.Client, abuseipdbClient *abuseipdb.Client, dnsClient *dns.Client) *API {
	server := echo.New()

	server.HideBanner = true
	server.HidePort = true

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	service := NewService(
		NewAddressStore(db),
		NewMetadataStore(db),
		cfClient,
		abuseipdbClient,
		dnsClient,
	)

	return &API{
		abuseipdbClient: abuseipdbClient,
		cfClient:        cfClient,
		Service:         service,
		Server:          server,
		Database:        db,
		dnsClient:       dnsClient,
	}
}

func (api *API) init() {
	// Blocklist Routes
	{
		api.Server.Add("GET", "/api/v1/addresses", api.handleAddressesGetAll)
		api.Server.Add("GET", "/api/v1/addresses/:ip", api.handleAddressesGetOne)
		api.Server.Add("DELETE", "/api/v1/addresses/:ip", api.handleAddressesDelete)
		api.Server.Add("POST", "/api/v1/addresses", api.handleAddressesPost)
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
