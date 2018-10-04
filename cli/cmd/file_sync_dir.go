package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(syncDirCmd)
	syncDirCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "recursively search for file")
}

var syncDirCmd = &cobra.Command{
	Use:   "syncd",
	Short: "sync a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		sr := strings.Split(root, "://")

		if len(sr) != 2 {
			log.Fatal("must include {driver}://{path} in argument eg. file:///home")
		}

		driver := sr[0]
		p := sr[1]

		var fs storage.Store
		var err error

		switch driver {
		case "dospace":
			fs, err = storage.NewDOSpace(
				viper.Get("do_space.access_key").(string),
				viper.Get("do_space.secret_key").(string),
				viper.Get("do_space.region").(string),
				viper.Get("do_space.bucket").(string),
				true,
			)
		default:
			fs, err = storage.NewFilesystem(p)
		}

		if err != nil {
			log.Fatal(err)
		}

		err = fs.Delete(context.Background(), p)
		if err != nil {
			log.Fatal(err)
		}
	},
}
