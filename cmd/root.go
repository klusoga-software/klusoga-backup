package cmd

import (
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "klusoga-backup",
}

var (
	Target    types.TargetType
	Username  string
	Password  string
	Host      string
	Port      int
	Databases []string
	Path      string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
