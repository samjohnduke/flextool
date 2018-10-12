package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"
	"sort"
	"time"

	"github.com/samjohnduke/flextool/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fileCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&listRecursive, "recursive", "R", false, "recursively search for file")
	listCmd.PersistentFlags().BoolVarP(&listAll, "all", "a", false, "show all files")
	listCmd.PersistentFlags().BoolVarP(&listAlmostAll, "almost-all", "A", false, "show all except . and ..")
	listCmd.PersistentFlags().BoolVarP(&listLong, "long", "l", false, "list files with complete details")
	listCmd.PersistentFlags().BoolVarP(&listWithSize, "with-size", "s", false, "show file size")
	listCmd.PersistentFlags().BoolVarP(&listHuman, "human", "H", false, "show numbers in human readable format eg 1.2Gb")
	listCmd.PersistentFlags().BoolVarP(&listSortBySize, "sort-size", "S", false, "sort files by size")
	listCmd.PersistentFlags().StringVar(&listSort, "sort", "", "sort by ${WORD}")
	listCmd.PersistentFlags().BoolVarP(&listSortByTime, "sort-time", "T", false, "sort by time")
	listCmd.PersistentFlags().BoolVarP(&listSortByName, "sort-name", "X", true, "sort by file name")
}

var listRecursive bool  // -R
var listAll bool        // -a
var listAlmostAll bool  // -A
var listLong bool       // -l
var listWithSize bool   // -s
var listHuman bool      // -h
var listSortBySize bool // -S
var listSort string     // sort=WORD
var listSortByTime bool // -T
var listSortByName bool // -X

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "show all the files at the specified URI",
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
			panic(err)
		}

		files, err := fs.List(context.Background(), name, storage.ListOpts{
			Recursive: listRecursive,
			All:       listAll,
			AlmostAll: listAlmostAll,
		})
		if err != nil {
			panic(err)
		}

		var fbs []storage.Blob
		for _, r := range files {
			fbs = append(fbs, r)
		}

		sortBlobsByName(fbs)

		for _, r := range fbs {
			fmt.Println(r.Name())
		}
	},
}

func sortBlobsByName(bs []storage.Blob) {
	sort.Slice(bs, func(i, j int) bool {
		return bs[i].Name() < bs[j].Name()
	})
}

func sortBlobsBySize(bs []storage.Blob) {
	sort.Slice(bs, func(i, j int) bool {
		var o, t int64
		s1, _ := bs[i].Stat()
		if s1 == nil {
			o = 0
		} else {
			o = s1.Size
		}

		s2, _ := bs[j].Stat()
		if s2 == nil {
			t = 0
		} else {
			t = s2.Size
		}

		return o < t
	})
}

func sortBlobsByTime(bs []storage.Blob) {
	sort.Slice(bs, func(i, j int) bool {
		var o, t time.Time
		s1, _ := bs[i].Stat()
		if s1 == nil {
			o = time.Now()
		} else {
			o = s1.LastModified
		}

		s2, _ := bs[j].Stat()
		if s2 == nil {
			t = time.Now()
		} else {
			t = s2.LastModified
		}

		return o.After(t)
	})
}
