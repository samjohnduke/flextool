package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pgCmd)
}

var pgCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Manage a postresql instance",
}
