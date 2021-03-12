package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hostinger/dnsrbl/pkg/alerters"
	"github.com/hostinger/dnsrbl/pkg/checkers"
	"github.com/hostinger/dnsrbl/pkg/endpoints"
	"github.com/hostinger/dnsrbl/pkg/hbl"
)

// @title Hostinger Block List API
// @version 1.0
// @description Hostinger HTTP service for managing IP address block lists.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := hbl.InitDB(ctx)
	if err != nil {
		log.Printf("Database: %s", err)
	}

	switch os.Getenv("ENVIRONMENT") {
	case "PRODUCTION":
		// Endpoints
		endpoints.Register(endpoints.NewCloudflareEndpoint())
		endpoints.Register(endpoints.NewPDNSEndpoint())
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(db))
		// Alerters
		alerters.Register(alerters.NewSlackAlerter())
	case "STAGING":
		// Endpoints
		endpoints.Register(endpoints.NewPDNSEndpoint())
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(db))
		// Alerters
		alerters.Register(alerters.NewSlackAlerter())
	default:
		// Endpoints
		endpoints.Register(endpoints.NewPDNSEndpoint())
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(db))
	}

	api := hbl.NewAPI(db)

	go func() {
		if err := api.Start(os.Getenv("HBL_LISTEN_ADDRESS"), os.Getenv("HBL_LISTEN_PORT")); err != nil {
			log.Fatalf("Failure: %s", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		api.Stop()
		os.Exit(0)
	}()

	select {}
}
