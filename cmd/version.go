package cmd

import (
	"fmt"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/build"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(build.Version)
	},
}
