package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hostinger/dnsrbl/pkg/abuseipdb"
	"github.com/hostinger/dnsrbl/pkg/cloudflare"
	"github.com/hostinger/dnsrbl/pkg/dns"
	"github.com/hostinger/dnsrbl/pkg/hbl"
)

// @title Hostinger Block List API
// @version 1.0
// @description Hostinger HTTP service for managing IP address block lists.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// Database
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	db, err := hbl.InitDB(ctx,
		os.Getenv("HBL_MYSQL_USERNAME"),
		os.Getenv("HBL_MYSQL_PASSWORD"),
		os.Getenv("HBL_MYSQL_HOST"),
		os.Getenv("HBL_MYSQL_PORT"),
		os.Getenv("HBL_MYSQL_DATABASE"))
	if err != nil {
		log.Printf("Failed to establish connection to the database: %s", err)
	}

	// Cloudflare
	cfClient, err := cloudflare.NewClient(os.Getenv("CF_API_ACCOUNT"), os.Getenv("CF_API_EMAIL"), os.Getenv("CF_API_KEY"))
	if err != nil {
		log.Printf("Failed to initialize Cloudflare client: %s", err)
	}

	// AbuseIPDB
	abuseipdbClient, err := abuseipdb.NewClient(os.Getenv("ABUSEIPDB_API_KEY"))
	if err != nil {
		log.Printf("Failed to initialize AbuseIPDB client: %s", err)
	}

	// PowerDNS
	dnsClient, err := dns.NewClient(os.Getenv("PDNS_API_SCHEME"),
		os.Getenv("PDNS_API_HOST"),
		os.Getenv("PDNS_API_PORT"),
		os.Getenv("PDNS_API_KEY"),
	)
	if err != nil {
		log.Printf("Failed to initialize PowerDNS client: %s", err)
	}

	api := hbl.NewAPI(db, cfClient, abuseipdbClient, dnsClient)

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
