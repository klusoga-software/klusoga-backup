package cmd

import (
	"errors"
	"fmt"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/destination"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"github.com/spf13/cobra"
	"log/slog"
)

var (
	Bucket string
	Region string
	Name   string
)

func init() {
	destinationAddAwsCmd.Flags().StringVarP(&Bucket, "bucket", "b", "", "The name of the bucket. You can uses slashes to add directories")
	destinationAddAwsCmd.Flags().StringVarP(&Region, "region", "r", "", "The aws region of the bucket")
	destinationAddCmd.PersistentFlags().StringVarP(&Name, "name", "n", "", "The name of the destination")

	destinationAddAwsCmd.MarkFlagRequired("bucket")
	destinationAddAwsCmd.MarkFlagRequired("region")
	destinationAddCmd.MarkFlagRequired("name")

	destinationAddCmd.AddCommand(destinationAddAwsCmd)
	destinationCmd.AddCommand(destinationListCmd, destinationAddCmd, destinationDeleteCmd)
}

func init() {
	rootCmd.AddCommand(destinationCmd)
}

var destinationCmd = &cobra.Command{
	Use: "destinations",
}

var destinationListCmd = &cobra.Command{
	Use: "list",
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

var destinationAddCmd = &cobra.Command{
	Use: "add",
}

var destinationAddAwsCmd = &cobra.Command{
	Use: "aws",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := destination.AddDestination(types.Destination{
			Name:   Name,
			Type:   types.Aws,
			Region: Region,
			Bucket: Bucket,
		})
		if err != nil {
			return err
		}

		slog.Info("Destination was added", "name", Name)
		return nil
	},
}

var destinationDeleteCmd = &cobra.Command{
	Use: "delete [name]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("need argument name")
		}

		err := destination.DeleteDestination(args[0])
		if err != nil {
			return err
		}

		slog.Info("Destination was deleted", "name", args[0])
		return nil
	},
}
