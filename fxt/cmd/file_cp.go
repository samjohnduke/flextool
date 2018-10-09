package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(copyCmd)
	copyCmd.PersistentFlags().BoolVarP(&copyBackup, "backup", "b", false, "backup duplicate files before copying")
	copyCmd.PersistentFlags().BoolVarP(&copyBackupVersioned, "version", "B", false, "backup files with a file version number")
	copyCmd.PersistentFlags().BoolVarP(&copyNoClobber, "clobber", "c", true, "overwrite files if they exist")
	copyCmd.PersistentFlags().BoolVarP(&copyInteractive, "interactive", "i", false, "copy files in interactive mode")
	copyCmd.PersistentFlags().BoolVarP(&copyRescursive, "recursive", "r", false, "copy files recursively")
}

var copyBackup bool          // -b
var copyBackupVersioned bool // -B
var copyNoClobber bool       // -n
var copyInteractive bool     // -i
var copyRescursive bool      // -r

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

		toStat, err := fromFS.Stat(context.TODO(), fromPath)
		var fromExists bool = true
		if err != nil {
			log.Println(err)
			if os.IsNotExist(err) {
				fromExists = false
			} else {
				log.Fatal(err)
			}
		}

		log.Println(toStat)

		_, err = toFS.Stat(context.TODO(), toPath)
		var toExists bool
		if err != nil {
			if os.IsNotExist(err) {
				toExists = false
			} else {
				log.Fatal(err)
			}
		}

		if !fromExists {
			fmt.Printf("file %s does not exist\n", fromPath)
			return
		}

		if toExists {
			// handle clobber / backup / versioning
		}

		fromBlob := fromFS.New(context.TODO(), "")
		if err != nil {
			log.Fatal(err)
		}

		if fromBlob.IsDir() && !copyRescursive {
			fmt.Printf("file %s is a directory. use -r to copy recursively\n", fromPath)
			return
		}

		if fromBlob.IsDir() {

		}

		log.Println(fromExists, toExists)

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
