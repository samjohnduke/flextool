package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pgCmd)
}

var pgCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "Backup a postresql instance",
	Long:  `Backup one or more databases in a postgresql instance`,
}
