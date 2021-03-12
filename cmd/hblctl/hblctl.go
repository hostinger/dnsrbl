package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hostinger/dnsrbl/sdk"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	hblScheme string
	hblHost   string
	hblPort   string
	hblKey    string
)

var client sdk.Client

var rootCmd = &cobra.Command{
	Use:   "hblctl",
	Short: "Hostinger Block List CLI",
	Long:  "Application which helps interact with Hostinger Block List API service.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if hblScheme == "" && os.Getenv("HBL_API_SCHEME") == "" {
			return errors.New("Flag or environment variable for 'HBL_API_SCHEME' is required. ")
		}
		if hblHost == "" && os.Getenv("HBL_API_HOST") == "" {
			return errors.New("Flag or environment variable for 'HBL_API_HOST' is required. ")
		}
		if hblPort == "" && os.Getenv("HBL_API_PORT") == "" {
			return errors.New("Flag or environment variable for 'HBL_API_PORT' is required. ")
		}
		if hblKey == "" && os.Getenv("HBL_API_KEY") == "" {
			return errors.New("Flag or environment variable for 'HBL_API_KEY' is required. ")
		}
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initClient)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.hblctl.yaml)")
	rootCmd.PersistentFlags().StringVar(
		&hblScheme, "hbl-api-scheme", os.Getenv("HBL_API_SCHEME"), "Scheme for connecting to the HBL API. (HBL_API_SCHEME)")
	rootCmd.PersistentFlags().StringVar(
		&hblHost, "hbl-api-host", os.Getenv("HBL_API_HOST"), "Host for connecting to the HBL API. (HBL_API_HOST)")
	rootCmd.PersistentFlags().StringVar(
		&hblPort, "hbl-api-port", os.Getenv("HBL_API_PORT"), "Port for connecting to the HBL API. (HBL_API_PORT)")
	rootCmd.PersistentFlags().StringVar(
		&hblKey, "hbl-api-key", os.Getenv("HBL_API_KEY"), "Key for connecting to the HBL API. (HBL_API_KEY)")
}

func initClient() {
	client = sdk.NewClient(hblKey, fmt.Sprintf("%s://%s:%s/api/v1", hblScheme, hblHost, hblPort))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".hblctl")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
