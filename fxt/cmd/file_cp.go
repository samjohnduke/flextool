package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path"

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
		// Load to and from paths as a URI
		from, err := url.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}

		to, err := url.Parse(args[1])
		if err != nil {
			log.Fatal(err)
		}

		// Get the path that we are working with
		fromPath, fromName := path.Split(from.Path)
		toPath, toName := path.Split(to.Path)

		// Load the Driver that we can use to copy with
		fromFS, err := getDriver(from.Scheme, fromPath)
		if err != nil {
			log.Fatal(err)
		}

		toFS, err := getDriver(to.Scheme, toPath)
		if err != nil {
			log.Fatal(err)
		}

		// Load up the blob for copying to and from
		fromBlob := fromFS.New(context.TODO(), fromName)
		toBlob := toFS.New(context.TODO(), toName)

		// Exit early if there is nothing to copy
		if !fromBlob.Exists() {
			fmt.Printf("file %s does not exist\n", fromPath)
			return
		}

		if toBlob.Exists() {
			// handle clobber / backup / versioning
		}

		// If the copy from blob is a directory but we aren't allowed to copy
		// exit early and print a short message to say how to copy the dir
		if fromBlob.IsDir() && !copyRescursive {
			fmt.Printf("file %s is a directory. use -r to copy recursively\n", fromPath)
			return
		}
		// Do the copy as we now now everything that need to start
		if fromBlob.IsDir() {

		} else {
			_, err := io.Copy(toBlob, fromBlob)
			if err != nil {
				panic(err)
			}

			toBlob.Close()
			fromBlob.Close()
		}
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
