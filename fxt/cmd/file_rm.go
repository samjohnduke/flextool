package cmd

import (
	"context"
	"log"
	"net/url"
	"path"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(removeCmd)
	removeCmd.PersistentFlags().BoolVarP(&removePreserveRoot, "preserve-root", "p", true, "ensure not to remove the root directory")
	removeCmd.PersistentFlags().BoolVarP(&removeRecursive, "recursive", "R", false, "delete folders recursively")
	removeCmd.PersistentFlags().BoolVarP(&removeEmptyDirectories, "empty", "d", false, "remove directory only if it is empty")
	removeCmd.PersistentFlags().BoolVarP(&removeInteractive, "interactive", "i", false, "interactive mode for deleting files")
}

var removePreserveRoot bool     // -p
var removeRecursive bool        // -R
var removeEmptyDirectories bool // -d
var removeInteractive bool      // -i

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "delete a file",
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

		err = fs.Delete(context.Background(), name)
		if err != nil {
			log.Fatal(err)
		}
	},
}
