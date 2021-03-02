package dnsrbl

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var server *http.Server

func Start(address string) {
	router := echo.New()
	router.Validator = NewValidator()

	router.GET("/health", HealthHandler)
	router.POST("/api/v1/blocklist", BlockHandler)
	router.GET("/api/v1/blocklist", GetAllHandler)
	router.GET("/api/v1/blocklist/:address", GetHandler)
	router.DELETE("/api/v1/blocklist/:address", UnblockHandler)

	server = &http.Server{
		Addr:    address,
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())
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
