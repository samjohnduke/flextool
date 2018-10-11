package cmd

import (
	"context"
	"fmt"
	"io"
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
		// Load to and from paths as a URI
		from, err := url.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}

		to, err := url.Parse(args[1])
		if err != nil {
			log.Fatal(err)
		}

		// If no driver is specified in the URI it means we are working with a
		// local filesytem path
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

		// Get the path that we are working with
		fromPath := from.Path
		toPath := to.Path

		// Load the Driver that we can use to copy with
		fromFS, err := getDriver(fromDriver, fromPath)
		if err != nil {
			log.Fatal(err)
		}

		toFS, err := getDriver(toDriver, toPath)
		if err != nil {
			log.Fatal(err)
		}

		// Determine if the copy from file/dir exists
		_, err = fromFS.Stat(context.TODO(), "")
		var fromExists bool = true
		if err != nil {
			if os.IsNotExist(err) {
				fromExists = false
			} else {
				log.Fatal(err)
			}
		}

		// Determine if the copy to file/dir exists
		_, err = toFS.Stat(context.TODO(), "")
		var toExists bool
		if err != nil {
			if os.IsNotExist(err) {
				toExists = false
			} else {

				log.Fatal(err)
			}
		}

		// Exit early if there is nothing to copy
		if !fromExists {
			fmt.Printf("file %s does not exist\n", fromPath)
			return
		}

		if toExists {
			// handle clobber / backup / versioning
		}

		// Load up the blob for copying from
		fromBlob := fromFS.New(context.TODO(), "")

		// If the copy from blob is a directory but we aren't allowed to copy
		// exit early and print a short message to say how to copy the dir
		if fromBlob.IsDir() && !copyRescursive {
			fmt.Printf("file %s is a directory. use -r to copy recursively\n", fromPath)
			return
		}
		// Do the copy as we now now everything that need to start
		if fromBlob.IsDir() {

		} else {
			toBlob := fromFS.New(context.TODO(), "")
			fromBlob.(*storage.File).EnsureReadable()
			toBlob.(*storage.File).EnsureWriteable()
			log.Println(fromBlob, toBlob)
			bytes, err := io.Copy(toBlob, fromBlob)
			if err != nil {
				panic(err)
			}

			toBlob.Close()
			fromBlob.Close()

			log.Println("copied bytes:", bytes)
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
