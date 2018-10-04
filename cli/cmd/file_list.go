package cmd

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "R", false, "recursively search for file")
}

var recursive bool  // -R
var all bool        // -a
var almostAll bool  // -A
var long bool       // -l
var withSize bool   // -s
var sortBySize bool // -S
var sort string     // sort=WORD
var sortByTime bool // -T
var sortByName bool // -X

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "show all the files in the application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]

		u, err := url.Parse(root)
		var driver string
		var p string

		if u.Scheme == "" {
			driver = "file"

		} else {
			driver = u.Scheme
		}
		p = u.Path

		base, name := path.Split(p)

		var fs storage.Store

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
