package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fileCmd)
}

var verbose bool

var fileCmd = &cobra.Command{
	Use:   "files",
	Short: "Use the files storage interface",
}
