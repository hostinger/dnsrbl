package main

import (
	"errors"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var blockCmd = &cobra.Command{
	Use:  "block <ip> <author> <comment>",
	Args: cobra.ExactArgs(3),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if net.ParseIP(args[0]) == nil {
			return errors.New("Argument 'IP' must be a valid IP address")
		}
		return nil
	},
	Short: "Block an IP address on Endpoints.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Block(cmd.Context(), args[0], args[1], args[2]); err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Print("Action executed successfully")
	},
}

func init() {
	rootCmd.AddCommand(blockCmd)
}
