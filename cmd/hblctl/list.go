package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"text/tabwriter"

	"github.com/hostinger/hbl/sdk"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:  "list [<ip>]",
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && net.ParseIP(args[0]) == nil {
			return errors.New("Argument 'IP' must be a valid IP address")
		}
		return nil
	},
	Short: "Get one or all addresses from database.",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
		if len(args) > 0 {
			address, err := client.GetOne(cmd.Context(), args[0])
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			writeAddressesHeader(w)
			writeAddressesTable(w, address)
			w.Flush()
			return
		}
		addresses, err := client.GetAll(cmd.Context())
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		writeAddressesHeader(w)
		for _, address := range addresses {
			writeAddressesTable(w, address)
		}
		w.Flush()
	},
}

func writeAddressesHeader(w io.Writer) {
	fmt.Fprint(w, "IP\tACTION\tAUTHOR\tCOMMENT\tCREATED_AT\n")
}

func writeAddressesTable(w io.Writer, args ...*sdk.Address) {
	for _, address := range args {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			address.IP, address.Action, address.Author, address.Comment, address.CreatedAt)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
