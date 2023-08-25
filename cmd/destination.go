package cmd

import (
	"fmt"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/destination"
	"github.com/spf13/cobra"
	"log/slog"
)

func init() {
	rootCmd.AddCommand(destinationCmd)
}

var destinationCmd = &cobra.Command{
	Use: "destinations",
	RunE: func(cmd *cobra.Command, args []string) error {
		destinations, err := destination.ListDestinations()
		if err != nil {
			return err
		}

		if len(destinations) == 0 {
			slog.Info("No destinations found")
			return nil
		}

		for _, d := range destinations {
			fmt.Println(d.String())
		}
		return nil
	},
}
