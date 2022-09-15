package cmd

import (
	"fmt"
	"github.com/s4kibs4mi/twilfe/config"
	"os"

	"github.com/spf13/cobra"
)

var (
	// RootCmd is the root command of nur service
	RootCmd = &cobra.Command{
		Use:   "twilfe",
		Short: "An ordering service for small cafe using Twilio & Shopemaa",
		Long:  `An ordering service for small cafe using Twilio & Shopemaa`,
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	if err := config.LoadConfig(); err != nil {
		fmt.Println("Failed to read config : ", err)
		os.Exit(1)
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
