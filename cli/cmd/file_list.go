package cmd

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "recursively search for file")
}

var recursive bool

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "show all the files in the application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		sr := strings.Split(root, "://")

		var driver string
		var p string

		if len(sr) == 1 {
			driver = "file"
			p = sr[0]
		} else {
			driver = sr[0]
			p = sr[1]
		}

		base, name := path.Split(p)

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
			fs, err = storage.NewFilesystem(base)
		}

		if err != nil {
			panic(err)
		}

		files, err := fs.List(context.Background(), name, storage.ListOpts{
			Recursive: recursive,
		})
		if err != nil {
			panic(err)
		}

		for _, r := range files {
			fmt.Println(path.Join(r.Name()))
		}
	},
}
