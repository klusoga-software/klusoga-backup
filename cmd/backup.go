package cmd

import (
	"github.com/klusoga-software/klusoga-backup-agent/pkg/backup"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"github.com/spf13/cobra"
	"log/slog"
)

func init() {
	backupCmd.Flags().StringVarP((*string)(&Target), "target", "t", "", "Target of what you want to backup")
	backupCmd.Flags().StringVarP(&Username, "username", "u", "", "Username")
	backupCmd.Flags().StringVarP(&Password, "password", "p", "", "Password")
	backupCmd.Flags().StringVar(&Host, "host", "", "Host of your database")
	backupCmd.Flags().IntVar(&Port, "port", 0, "Port of your database")
	backupCmd.Flags().StringSliceVar(&Databases, "databases", []string{}, "The list of databases you want to backup")
	backupCmd.Flags().StringVar(&Path, "path", "", "The path to store the backup")

	backupCmd.MarkFlagRequired("target")
	backupCmd.MarkFlagRequired("username")
	backupCmd.MarkFlagRequired("password")
	backupCmd.MarkFlagRequired("host")
	backupCmd.MarkFlagRequired("port")
	backupCmd.MarkFlagRequired("path")
}

func init() {
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use: "backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		var bak backup.Target

		switch Target {
		case types.TargetTypeMssql:
			slog.Info("Run Mssql Backup")
			mssqlBackup := backup.NewMssqlTarget(backup.MssqlParameters{
				Host:      Host,
				Port:      Port,
				Username:  Username,
				Password:  Password,
				Databases: Databases,
				Path:      Path,
			})

			bak = mssqlBackup
		default:
			slog.Error("Target Type not found")
			return nil
		}

		if err := bak.Backup(); err != nil {
			return err
		}

		return nil
	},
}
