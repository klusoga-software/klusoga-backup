package cmd

import (
	"errors"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/backup"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/destination"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"log/slog"
	"os"
)

func init() {
	backupCmd.Flags().StringVarP((*string)(&Target), "target", "t", "", "Target of what you want to backup")
	backupCmd.Flags().StringVarP(&Username, "username", "u", "", "Username")
	backupCmd.Flags().StringVarP(&Password, "password", "p", "", "Password")
	backupCmd.Flags().StringVar(&Host, "host", "", "Host of your database")
	backupCmd.Flags().IntVar(&Port, "port", 0, "Port of your database")
	backupCmd.Flags().StringSliceVar(&Databases, "databases", []string{}, "The list of databases you want to backup")
	backupCmd.Flags().StringVar(&Path, "path", "", "The path to store the backup")
	backupCmd.Flags().StringVar(&Destination, "destination", "", "Where you want to store your backup")

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

		files, err := bak.Backup()
		if err != nil {
			return err
		}

		if Destination != "" {
			slog.Info("Destination specified upload backup")
			var destinationFile types.DestinationFile
			var dest types.Destination
			file, err := os.Open(os.Getenv("DESTINATION_FILE_PATH"))
			if err != nil {
				return err
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			yaml.Unmarshal(data, &destinationFile)

			for _, d := range destinationFile.Destinations {
				if d.Name == Destination {
					dest = d
				}
			}

			switch dest.Type {
			case types.Aws:
				des = destination.NewS3BucketDestination(destination.S3DestinationParams{
					Bucket: dest.Bucket,
					Region: dest.Region,
				})
			default:
				return errors.New("destination Type don't match")
			}

			err = des.UploadFiles(files, "klusoga")
			if err != nil {
				return err
			}
		}

		return nil
	},
}
