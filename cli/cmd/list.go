package cmd

import (
	"context"
	"log"
	"path"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "file", "Driver to use as store")
	listCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "recursively search for file")
}

var driver string
var recursive bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show all the files in the application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		fs, err := storage.NewFilesystem(root)
		if err != nil {
			panic(err)
		}

		files, err := fs.List(context.Background(), "", storage.ListOpts{
			Recursive: recursive,
		})
		if err != nil {
			panic(err)
		}

		for _, r := range files {
			log.Println(path.Join(r.Path(), r.Name()))
		}
	},
}
