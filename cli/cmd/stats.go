package cmd

import (
	"log"

	"github.com/shirou/gopsutil/load"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get statistics about the current state of the machine",
	Run: func(cmd *cobra.Command, args []string) {
		stat, err := load.Avg()
		if err != nil {
			panic(err)
		}

		log.Println(stat)
	},
}
