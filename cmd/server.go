package cmd

import (
	"github.com/cs5224virgo/virgo/db"
	"github.com/cs5224virgo/virgo/internal/api"
	"github.com/cs5224virgo/virgo/internal/datalayer"
	"github.com/cs5224virgo/virgo/internal/socket"
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
		dbname := viper.GetString("db.name")
		logger.Info("db.name= ", dbname)

		db, err := db.InitDB()
		if err != nil {
			logger.Fatal("cannot connect to DB:", err)
		}
		logger.Info("Connected to DB successfully")

		// init datalayer
		data := datalayer.NewDataLayer(db)

		// init websocket
		hub := socket.NewWebSocketHub(data)
		go hub.Run()

		// init apiserver
		api := api.NewAPIServer(data, hub)

		// lmao
		logger.Info(`
   _            .                                         
  u            @88>                                       
 88Nu.   u.    %8P      .u    .                      u.   
'88888.o888c    .     .d88B :@8c       uL      ...ue888b  
 ^8888  8888  .@88u  ="8888f8888r  .ue888Nc..  888R Y888r 
  8888  8888 ''888E'   4888>'88"  d88E'"888E'  888R I888> 
  8888  8888   888E    4888> '    888E  888E   888R I888> 
  8888  8888   888E    4888>      888E  888E   888R I888> 
 .8888b.888P   888E   .d888L .+   888E  888E  u8888cJ888  
  ^Y8888*""    888&   ^"8888*"    888& .888E   "*888*P"   
    'Y"        R888"     "Y"      *888" 888&     'Y"      
                ""                 '"   "888E             
                                  .dWi   '88E             
                                  4888~  J8%              
                                   ^"===*"'               

		`)

		// run
		api.Run()
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
