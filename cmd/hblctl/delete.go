package main

import (
	"errors"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:  "delete <ip>",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if net.ParseIP(args[0]) == nil {
			return errors.New("Argument 'IP' must be a valid IP address")
		}
		return nil
	},
	Short: "Delete an IP address on Endpoints.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Delete(cmd.Context(), args[0]); err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Print("Action executed successfully")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
