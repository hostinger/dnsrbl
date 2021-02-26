package dnsrbl

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hostinger/dnsrbl/database"
	"github.com/labstack/echo/v4"
)

var server *http.Server

// Start ...
func Start(address string) {
	handler := NewHandler(
		NewMySQLAddressStore(database.DB),
	)

	router := echo.New()
	router.POST("/api/v1/block", handler.BlockHandler)
	router.GET("/api/v1/search/:address", handler.SearchHandler)
	router.DELETE("/api/v1/unblock/:address", handler.UnblockHandler)

	server = &http.Server{
		Addr:    address,
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())
}

// Stop ...
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
