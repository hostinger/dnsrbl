package dnsrbl

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/hostinger/dnsrbl/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var server *echo.Echo

func Start(address string) {
	server = echo.New()
	server.Validator = NewValidator()

	server.GET("/health", HealthHandler)
	server.POST("/api/v1/blocklist", BlockHandler)
	server.GET("/api/v1/blocklist", GetAllHandler)
	server.GET("/api/v1/blocklist/:address", GetHandler)
	server.DELETE("/api/v1/blocklist/:address", UnblockHandler)

	server.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := server.Start(address); err != nil && err != http.ErrServerClosed {
		server.Logger.Fatal("Exiting...")
	}
}

func Stop() {
	log.Print("Shutting down server with timeout of 5 seconds...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to cleanly shutdown server: ", err)
	}
}
