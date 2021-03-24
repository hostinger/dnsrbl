package main

import (
	"errors"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:  "sync [<ip>]",
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && net.ParseIP(args[0]) == nil {
			return errors.New("Argument 'IP' must be a valid IP address")
		}
		return nil
	},
	Short: "Sync one or all addresses from database.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if err := client.SyncOne(cmd.Context(), args[0]); err != nil {
				log.Fatalf("Error: %s", err)
			}
			log.Print("Action executed successfully")
			return
		}
		if err := client.SyncAll(cmd.Context()); err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Print("Action executed successfully")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
