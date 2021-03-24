package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hostinger/hbl/pkg/alerters"
	"github.com/hostinger/hbl/pkg/checkers"
	"github.com/hostinger/hbl/pkg/endpoints"
	"github.com/hostinger/hbl/pkg/hbl"
	"go.uber.org/zap"
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

	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Logger: %s", err)
	}

	api := hbl.NewAPI(
		&hbl.Config{
			DB:     db,
			Host:   os.Getenv("HBL_LISTEN_ADDRESS"),
			Port:   os.Getenv("HBL_LISTEN_PORT"),
			Logger: logger,
		},
	)

	switch os.Getenv("ENVIRONMENT") {
	case "PRODUCTION":
		// Endpoints
		endpoints.Register(endpoints.NewCloudflareEndpoint(logger))
		endpoints.Register(endpoints.NewPDNSEndpoint(logger))
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(logger, db))
		// Alerters
		alerters.Register(alerters.NewSlackAlerter(logger))
	case "STAGING":
		// Endpoints
		endpoints.Register(endpoints.NewPDNSEndpoint(logger))
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(logger, db))
		// Alerters
		alerters.Register(alerters.NewSlackAlerter(logger))
	default:
		// Endpoints
		endpoints.Register(endpoints.NewPDNSEndpoint(logger))
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(logger, db))
	}

	go func() {
		if err := api.Start(); err != nil {
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
