package main

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/sm"
	"github.com/prologic/sm/client"
)

// createCmd represents the run command
var createCmd = &cobra.Command{
	Use:     "create [flags] <title>",
	Aliases: []string{"new"},
	Short:   "Creates a new sev with the given title",
	Long:    `This creates a new sev with the given title and optional level.`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uri := viper.GetString("uri")
		client := client.NewClient(uri, nil)

		level, err := cmd.Flags().GetInt("level")
		if err != nil {
			log.Errorf("error getting -l/--level flag: %s", err)
			os.Exit(1)
		}

		quiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			log.Errorf("error getting -q/--quiet flag: %s", err)
			os.Exit(1)
		}

		os.Exit(create(client, args[0], level, quiet))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().IntP(
		"level", "l", sm.DefaultSEVLevel,
		"Set SEV level",
	)

	createCmd.Flags().BoolP(
		"quiet", "q", false,
		"Only display numeric IDs",
	)
}

func create(client *client.Client, title string, level int, quiet bool) int {
	res, err := client.Create(title, level)
	if err != nil {
		log.Errorf("error creating event %s: %s", title, err)
		return 1
	}

	if quiet {
		fmt.Print(res[0].ID)
	} else {
		out, err := json.Marshal(res)
		if err != nil {
			log.Errorf("error encoding event result: %s", err)
			return 1
		}
		fmt.Printf(string(out))
	}

	return 0
}
