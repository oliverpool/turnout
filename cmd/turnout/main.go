package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
)

func main() {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		configDir = "."
	}
	configDir = filepath.Join(configDir, "turnout")

	rootCmd := &cobra.Command{
		Use:   "turnout",
		Short: "Port-sharing proxy for local development",
	}
	rootCmd.PersistentFlags().StringP("configdir", "c", configDir, "dir to store the config")

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the proxy on the given port",
		RunE:  serve,
	}
	serveCmd.Flags().StringP("address", "a", "localhost:9999", "IP to listen on")

	addCmd := &cobra.Command{
		Use:   "add [hostname] [config]",
		Short: "Add a proxy for a given hostname to the config",
		Long:  `Add a proxy for a given hostname to the config. Use an Int to specify a port to forward to.`,
		Args:  cobra.ExactArgs(2),
		RunE:  add,
	}

	rootCmd.AddCommand(serveCmd, addCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
