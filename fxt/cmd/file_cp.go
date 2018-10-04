package cmd

import (
	"log"
	"strings"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(copyCmd)
}

var copyCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy a file from 1 location to another",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		from := strings.Split(args[0], "://")

		if len(from) != 2 {
			log.Fatal("must include {driver}://{path} in argument eg. file:///home")
		}

		fromDriver := from[0]
		fromPath := from[1]

		var fromFS storage.Store
		var err error

		switch fromDriver {
		case "dospace":
			fromFS, err = storage.NewDOSpace(
				viper.Get("do_space.access_key").(string),
				viper.Get("do_space.secret_key").(string),
				viper.Get("do_space.region").(string),
				viper.Get("do_space.bucket").(string),
				true,
			)
		default:
			fromFS, err = storage.NewFilesystem(fromPath)
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println(fromFS)
	},
}
