package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hostinger/hbl/pkg/alerters"
	"github.com/hostinger/hbl/pkg/checkers"
	"github.com/hostinger/hbl/pkg/database"
	"github.com/hostinger/hbl/pkg/endpoints"
	"github.com/hostinger/hbl/pkg/hbl"
	"github.com/hostinger/hbl/pkg/logger"
	"go.uber.org/zap"
)

// @title Hostinger Block List API
// @version 1.0
// @description Hostinger HTTP service for managing IP address block lists.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	l := logger.NewLoggerFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.Init(ctx,
		os.Getenv("HBL_MYSQL_USERNAME"),
		os.Getenv("HBL_MYSQL_PASSWORD"),
		os.Getenv("HBL_MYSQL_HOST"),
		os.Getenv("HBL_MYSQL_PORT"),
		os.Getenv("HBL_MYSQL_DATABASE"),
	)
	if err != nil {
		l.Fatal(
			"Failed to initialize database",
			zap.String("host", os.Getenv("HBL_MYSQL_HOST")),
			zap.String("port", os.Getenv("HBL_MYSQL_PORT")),
			zap.Error(err),
		)
	}

	r := hbl.NewMySQLRepository(l, db)
	s := hbl.NewDefaultService(l, r)
	h := hbl.NewDefaultHandler(l, s)

	api := hbl.NewAPI(
		&hbl.Config{
			Handler: h,
			Logger:  l,
			Host:    os.Getenv("HBL_LISTEN_ADDRESS"),
			Port:    os.Getenv("HBL_LISTEN_PORT"),
		},
	)

	api.Initialize(ctx)

	switch os.Getenv("ENVIRONMENT") {
	case "PRODUCTION":
		// Endpoints
		endpoints.Register(endpoints.NewCloudflareEndpoint(l))
		endpoints.Register(endpoints.NewPDNSEndpoint(l))
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(l, db))
		// Alerters
		alerters.Register(alerters.NewSlackAlerter(l))
	case "STAGING":
		// Endpoints
		endpoints.Register(endpoints.NewPDNSEndpoint(l))
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(l, db))
		// Alerters
		alerters.Register(alerters.NewSlackAlerter(l))
	default:
		// Endpoints
		endpoints.Register(endpoints.NewPDNSEndpoint(l))
		// Checkers
		checkers.Register(checkers.NewAbuseIPDBChecker(l, db))
	}

	go func() {
		api.Start()
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
