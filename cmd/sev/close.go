package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/sm/client"
)

// closeCmd represents the run command
var closeCmd = &cobra.Command{
	Use:     "close [flags] <id>",
	Aliases: []string{"done"},
	Short:   "Close the specified event",
	Long:    `This closes the specified event`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uri := viper.GetString("uri")
		client := client.NewClient(uri, nil)

		id := args[0]

		os.Exit(close(client, id))
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)
}

func close(client *client.Client, id string) int {
	err := client.Close(id)
	if err != nil {
		log.Errorf("error writing to event #%s: %s", id, err)
		return 1
	}

	return 0
}
