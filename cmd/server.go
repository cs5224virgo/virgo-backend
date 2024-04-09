package cmd

import (
	"github.com/cs5224virgo/virgo/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the virgo web server backend",
	Long:  `Runs the virgo web server backend`,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("server called")
		dbname := viper.GetString("dbname")
		logger.Info("dbname= ", dbname)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
