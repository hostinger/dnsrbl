package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hostinger/dnsrbl"
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
	cfgFile = flag.String("config.file", "config.yml", "Path to the configuration file.")
)

// @title Hostinger Block List API
// @version 1.0
// @description Hostinger HTTP service for managing IP address block lists.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	flag.Parse()

	db, err := dnsrbl.InitDB(*mysqlUsername, *mysqlPassword, *mysqlHost, *mysqlPort, *mysqlDatabase)
	if err != nil {
		log.Fatal("Failed to initialize connection to the database: ", err)
	}

	config, err := dnsrbl.NewConfigFromFile(*cfgFile)
	if err != nil {
		log.Fatalf("Failed to load configuration file: %s", err)
	}
	if err := config.Validate(); err != nil {
		log.Fatalf("Failed to validate configuration file: %s", err)
	}

	api := dnsrbl.NewAPI(config, db)

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
