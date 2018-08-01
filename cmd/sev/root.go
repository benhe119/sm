package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/sm"
)

var configFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "sev",
	Version: sm.FullVersion(),
	Short:   "Command-line client for sm",
	Long: `This is the command-line client for the sev manager daemon sm.

This lets you create, search, comment and manipulate the state of events.

This is the reference implementation of using the sm client library for
creating and mutating sev(s) (site events)`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set logging level
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
}

// Execute adds all child commands to the root command
// and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(
		&configFile, "config", "",
		"config file (default is $HOME/.sm.yaml)",
	)

	RootCmd.PersistentFlags().BoolP(
		"debug", "d", false,
		"Enable debug logging",
	)

	RootCmd.PersistentFlags().StringP(
		"uri", "u", "http://localhost:8000",
		"URI to connect to sm",
	)

	viper.BindPFlag("uri", RootCmd.PersistentFlags().Lookup("uri"))
	viper.SetDefault("uri", "http://localhost:8000/")

	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("debug", false)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".sm.yaml")
	}

	// from the environment
	viper.SetEnvPrefix("SM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
