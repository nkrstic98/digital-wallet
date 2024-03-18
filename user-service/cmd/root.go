package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"user-service/config"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "user-service",
	Short: "User microservice",
	Long:  "Microservice for user management",
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
