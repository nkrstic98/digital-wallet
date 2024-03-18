package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"transaction-service/config"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "transaction-service",
	Short: "Transaction microservice",
	Long:  "Microservice for handling user transactions and storing user balances.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cfg = &config.DefaultConfig
}
