package cmd

import (
	"github.com/robrotheram/gogallery/backend/api"
	"github.com/robrotheram/gogallery/backend/config"
	"github.com/robrotheram/gogallery/backend/datastore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCMD)
}

var serveCMD = &cobra.Command{
	Use:   "serve",
	Short: "Serve static site",
	Long:  "Serve static site",
	Run: func(cmd *cobra.Command, args []string) {

		config := config.LoadConfig()
		db := datastore.Open(config.Gallery.Basepath)
		defer db.Close()
		db.ScanPath(config.Gallery.Basepath)
		if len(args) == 1 {
			config.Server.Port = ":" + args[0]
		}
		api.NewGoGalleryAPI(config, db).Serve()
	},
}
