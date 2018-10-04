package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "absolute path to config")

	rootCmd.PersistentFlags().StringVar(&doAK, "do_access_key", "", "Access key for digital ocean spaces")
	rootCmd.PersistentFlags().StringVar(&doSK, "do_secret_key", "", "Secret key for digital ocean spaces")
	rootCmd.PersistentFlags().StringVar(&doRegion, "do_region", "", "Region for digital ocean spaces")
	rootCmd.PersistentFlags().StringVar(&doBucket, "do_bucket", "", "Bucket for digital ocean spaces")

	if !rootCmd.PersistentFlags().Changed("do_access_key") {
		viper.Set("do_spaces.access_key", doAK)
	}

	if !rootCmd.PersistentFlags().Changed("do_secret_key") {
		viper.Set("do_spaces.secret_key", doSK)
	}

	if !rootCmd.PersistentFlags().Changed("do_region") {
		viper.Set("do_spaces.region", doRegion)
	}

	if !rootCmd.PersistentFlags().Changed("do_bucket") {
		viper.Set("do_spaces.bucket", doBucket)
	}
}

var cfgFile string

var doAK string
var doSK string
var doRegion string
var doBucket string

var rootCmd = &cobra.Command{
	Use:   "fxt",
	Short: "",
	Long:  ``,
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

// Execute runs the root command from which all commands are connected
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
