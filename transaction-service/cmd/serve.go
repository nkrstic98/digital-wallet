package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
	"transaction-service/app"
	"transaction-service/db"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts user-service microservice",
	Long:  `Starts user-service microservice`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := app.Build(*cfg)
		if err != nil {
			panic(err)
		}
		defer db.CloseConnection()

		addr := fmt.Sprintf("%v:%v", cfg.Web.Host, cfg.Web.Port)
		slog.Info(fmt.Sprintf("Starting http server on port %v", cfg.Web.Port))
		go func() {
			if err = app.Run(addr); err != nil {
				panic(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		slog.Info("Shutting down server...")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
