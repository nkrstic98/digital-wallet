package cmd

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"transaction-service/config"
	"transaction-service/db"
)

var reinitCmd = &cobra.Command{
	Use:   "reinit",
	Short: "Recreates underlying database schema",
	Long:  "Recreates underlying database schema",
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Runs command on the database resource",
	Long:  "Runs command on the database resource",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := db.OpenConnection(lo.FromPtrOr(cfg, config.DefaultConfig))
		if err != nil {
			panic(err)
		}

		err = db.ReinitDatabase()
		if err != nil {
			panic(err)
		}

		err = db.CloseConnection()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(reinitCmd)
	reinitCmd.AddCommand(dbCmd)
}
