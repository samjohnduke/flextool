package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(statCmd)
	statCmd.PersistentFlags().BoolVarP(&statHuman, "human", "h", false, "print using human readable sizes eg. 1.2Gb")
}

var statHuman bool // -h

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "get stat information for a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		u, err := url.Parse(root)
		if err != nil {
			log.Fatal(err)
		}

		base, name := path.Split(u.Path)

		var fs storage.Store

		switch u.Scheme {
		case "dospace":
			fs, err = storage.NewDOSpace(
				viper.Get("do_space.access_key").(string),
				viper.Get("do_space.secret_key").(string),
				viper.Get("do_space.region").(string),
				viper.Get("do_space.bucket").(string),
				true,
			)
		default:
			fs, err = storage.NewFilesystem(base)
		}

		if err != nil {
			log.Fatal(err)
		}

		stat, err := fs.Stat(context.Background(), name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", stat)
	},
}
