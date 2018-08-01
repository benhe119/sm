package main

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/sm/client"
)

// infoCmd represents the run command
var infoCmd = &cobra.Command{
	Use:     "info [flags] <id>",
	Aliases: []string{"view"},
	Short:   "View information about an event",
	Long:    `This retrieves and views information about an event`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uri := viper.GetString("uri")
		client := client.NewClient(uri, nil)

		id := args[0]

		os.Exit(info(client, id))
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}

func info(c *client.Client, id string) int {
	res, err := c.GetEventByID(id)
	if err != nil {
		log.Errorf("error retrieving information for event #%s: %s", id, err)
		return 1
	}

	out, err := json.Marshal(res)
	if err != nil {
		log.Errorf("error encoding event: %s", err)
		return 1
	}

	fmt.Print(string(out))

	return 0
}
