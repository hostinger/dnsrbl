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
)

type API struct {
	Server  *echo.Echo
	Service Service
}

func NewAPI(db *sql.DB) *API {
	server := echo.New()
	server.HidePort = true
	server.HideBanner = true

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/swagger/*", echoSwagger.WrapHandler)

	service := NewService(NewMySQLRepository(db))

	return &API{
		Server:  server,
		Service: service,
	}
}

func (api *API) Start(host, port string) error {
	api.SetupRoutes()
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
		log.Printf("Failed clean shutdown, exiting with errors.")
	}
}
