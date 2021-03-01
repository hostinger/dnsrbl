package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hostinger/dnsrbl"
	"github.com/hostinger/dnsrbl/database"
)

var (
	mysqlUsername = os.Getenv("HBL_MYSQL_USERNAME")
	mysqlPassword = os.Getenv("HBL_MYSQL_PASSWORD")
	mysqlDatabase = os.Getenv("HBL_MYSQL_DATABASE")
	mysqlHost     = os.Getenv("HBL_MYSQL_HOST")
	mysqlPort     = os.Getenv("HBL_MYSQL_PORT")
)

var (
	listenAddress = os.Getenv("HBL_LISTEN_ADDRESS")
	listenPort    = os.Getenv("HBL_LISTEN_PORT")
)

// @title Hostinger Block List API
// @version 1.0
// @description Hostinger HTTP service for managing IP address block lists.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	if err := database.Init(mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase); err != nil {
		log.Fatal("Failed to initialize connection to the database: ", err)
	}

	go func() {
		dnsrbl.Start(fmt.Sprintf("%s:%s", listenAddress, listenPort))
	}()

	log.Printf("Listening on %s:%s", listenAddress, listenPort)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println()
		dnsrbl.Stop()
		database.DB.Close()
		log.Print("Exiting...")
		os.Exit(0)
	}()

	select {}
}
