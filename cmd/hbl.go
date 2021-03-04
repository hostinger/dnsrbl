package main

import (
	"context"
	"flag"
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

var (
	mysqlUsername = flag.String("db.username", os.Getenv("HBL_MYSQL_USERNAME"), "Username for database connection.")
	mysqlPassword = flag.String("db.password", os.Getenv("HBL_MYSQL_PASSWORD"), "Password for database connection.")
	mysqlDatabase = flag.String("db.database", os.Getenv("HBL_MYSQL_DATABASE"), "Name of the database.")
	mysqlHost     = flag.String("db.host", os.Getenv("HBL_MYSQL_HOST"), "Host for database connection.")
	mysqlPort     = flag.String("db.port", os.Getenv("HBL_MYSQL_PORT"), "Port for database connection.")
)

var (
	listenAddress = flag.String("listen.address", os.Getenv("HBL_LISTEN_ADDRESS"), "Listen address for HTTP server.")
	listenPort    = flag.String("listen.port", os.Getenv("HBL_LISTEN_PORT"), "Listen port for HTTP server.")
)

var (
	cloudflareKey     = flag.String("cf.api-key", os.Getenv("CF_API_KEY"), "Cloudflare API Key.")
	cloudflareEmail   = flag.String("cf.api-email", os.Getenv("CF_API_EMAIL"), "Cloudflare API Email.")
	cloudflareAccount = flag.String("cf.api-account", os.Getenv("CF_API_ACCOUNT"), "Cloudflare API Account.")
)

var (
	abuseipdbKey = flag.String("abuseipdb.key", os.Getenv("ABUSEIPDB_API_KEY"), "AbuseIPDB API Key.")
)

var (
	pdnsScheme = flag.String("pdns.scheme", os.Getenv("PDNS_API_SCHEME"), "PowerDNS API Scheme.")
	pdnsHost   = flag.String("pdns.host", os.Getenv("PDNS_API_HOST"), "PowerDNS API Host.")
	pdnsPort   = flag.String("pdns.port", os.Getenv("PDNS_API_PORT"), "PowerDNS API Port.")
	pdnsKey    = flag.String("pdns.key", os.Getenv("PDNS_API_KEY"), "PowerDNS API Key.")
)

var (
	cfgFile = flag.String("config.file", "config.yml", "Path to the configuration file.")
)

// @title Hostinger Block List API
// @version 1.0
// @description Hostinger HTTP service for managing IP address block lists.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	flag.Parse()

	// Database
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	db, err := hbl.InitDB(ctx, *mysqlUsername, *mysqlPassword, *mysqlHost, *mysqlPort, *mysqlDatabase)
	if err != nil {
		log.Fatalf("Failed to establish connection to the database: %s", err)
	}

	// Cloudflare
	cfClient, err := cloudflare.NewClient(*cloudflareAccount, *cloudflareEmail, *cloudflareKey)
	if err != nil {
		log.Printf("Failed to initialize Cloudflare client: %s", err)
	}

	// AbuseIPDB
	abuseipdbClient, err := abuseipdb.NewClient(*abuseipdbKey)
	if err != nil {
		log.Printf("Failed to initialize AbuseIPDB client: %s", err)
	}

	// PowerDNS
	dnsClient, err := dns.NewClient(*pdnsScheme, *pdnsHost, *pdnsPort, *pdnsKey)
	if err != nil {
		log.Printf("Failed to initalize PowerDNS client: %s", err)
	}

	api := hbl.NewAPI(db, cfClient, abuseipdbClient, dnsClient)

	go func() {
		api.Start(*listenAddress, *listenPort)
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
