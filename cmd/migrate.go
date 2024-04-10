package cmd

import (
	"github.com/cs5224virgo/virgo/db"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "perform migrations on the db",
	Long:  `Perform migrations on the db`,

	Run: func(cmd *cobra.Command, args []string) {
		err := db.Migrate()
		if err != nil {
			logger.Fatal("unable to perform migration:", err)
		}
		logger.Info("Migrations completed successfully")
		err = db.PrintVersion()
		if err != nil {
			logger.Fatal("unable to print migration version:", err)
		}
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current migration version",
	Long:  `Print current migration version`,
	Run: func(cmd *cobra.Command, args []string) {
		err := db.PrintVersion()
		if err != nil {
			logger.Fatal("unable to print migration version:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateVersionCmd)

}
