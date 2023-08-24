package cmd

import (
	"github.com/klusoga-software/klusoga-backup-agent/pkg/backup"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/destination"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/signal"
)

func init() {
	backupCmd.Flags().StringVarP((*string)(&Target), "target", "t", "", "Target of what you want to backup like mssql")
	backupCmd.Flags().StringVarP(&Username, "username", "u", "", "Username")
	backupCmd.Flags().StringVarP(&Password, "password", "p", "", "Password")
	backupCmd.Flags().StringVar(&Host, "host", "", "Host of your database")
	backupCmd.Flags().IntVar(&Port, "port", 0, "Port of your database")
	backupCmd.Flags().StringSliceVar(&Databases, "databases", []string{}, "The comma separated list of databases you want to backup")
	backupCmd.Flags().StringVar(&Path, "path", "", "The path to store the backup")
	backupCmd.Flags().StringVar(&Destination, "destination", "", "Where you want to store your backup")
	backupCmd.Flags().StringVar(&Schedule, "schedule", "", "The cron job schedule you want to run backups")

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
		var des destination.Destination

		switch Target {
		case types.TargetTypeMssql:
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

		if Schedule == "" {
			files, err := bak.Backup()
			if err != nil {
				return err
			}

			if Destination != "" {
				slog.Info("Destination specified upload backup")

				des, err = destination.GetDestinationByName(Destination)
				if err != nil {
					return err
				}

				err = des.UploadFiles(files, "klusoga")
				if err != nil {
					return err
				}
			}
		} else {
			slog.Info("Schedule specified. Run cronjob")

			forever := make(chan os.Signal, 1)
			signal.Notify(forever, os.Interrupt)

			c := cron.New(cron.WithSeconds())

			_, err := c.AddFunc(Schedule, func() {
				files, err := bak.Backup()
				if err != nil {
					slog.Error("Error while backup database", "error", err.Error())
					return
				}

				if Destination != "" {
					slog.Info("Destination specified upload backup")

					des, err = destination.GetDestinationByName(Destination)
					if err != nil {
						slog.Error("Error while receive destination", "error", err.Error())
						return
					}

					err = des.UploadFiles(files, "klusoga")
					if err != nil {
						slog.Error("Error while upload file", "error", err.Error())
						return
					}
				}
			})
			if err != nil {
				return err
			}

			c.Start()

			<-forever
			slog.Info("Stop schedule")
		}

		return nil
	},
}
