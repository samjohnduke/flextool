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
	fileCmd.AddCommand(removeCmd)
	removeCmd.PersistentFlags().BoolVarP(&removePreserveRoot, "preserve-root", "-p", true, "ensure not to remove the root directory")
	removeCmd.PersistentFlags().BoolVarP(&removeRecursive, "recursive", "-R", false, "delete folders recursively")
	removeCmd.PersistentFlags().BoolVarP(&removeEmptyDirectories, "empty", "-d", false, "remove directory only if it is empty")
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
