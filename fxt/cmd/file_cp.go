package cmd

import (
	"log"
	"net/url"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(copyCmd)
}

var copyBackup bool          // -b
var copyBackupVersioned bool // -B
var copyNoClobber bool       // -n
var copyInteractive bool     // -i

var copyCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy a file from 1 location to another",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		from, err := url.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}

		to, err := url.Parse(args[1])
		if err != nil {
			log.Fatal(err)
		}

		var fromDriver string
		if from.Scheme == "" {
			fromDriver = "file"
		} else {
			fromDriver = from.Scheme
		}

		var toDriver string
		if to.Scheme == "" {
			toDriver = "file"
		} else {
			toDriver = from.Scheme
		}

		fromPath := from.Path
		toPath := to.Path

		fromFS, err := getDriver(fromDriver, fromPath)
		if err != nil {
			log.Fatal(err)
		}

		toFS, err := getDriver(toDriver, toPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(fromFS, toFS)
	},
}

func getDriver(driver string, path string) (storage.Store, error) {
	switch driver {
	case "dospace":
		return storage.NewDOSpace(
			viper.Get("do_space.access_key").(string),
			viper.Get("do_space.secret_key").(string),
			viper.Get("do_space.region").(string),
			viper.Get("do_space.bucket").(string),
			true,
		)
	default:
		return storage.NewFilesystem(path)
	}
}
