package cmd

import (
	"log"

	"github.com/samjohnduke/flextool/archive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	pgCmd.AddCommand(pgbkupCmd)
}

var pgbkupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  ``,
	Run: func(c *cobra.Command, args []string) {
		pg := archive.PostgreSQL{
			Username: viper.Get("PG.username").(string),
			Path:     viper.Get("PG.path").(string),
			All:      viper.Get("PG.all").(bool),
			DB:       viper.Get("PG.db").(string),
		}

		archive, err := pg.Archive()
		if err != nil {
			panic(err)
		}

		log.Println(archive)
	},
}
