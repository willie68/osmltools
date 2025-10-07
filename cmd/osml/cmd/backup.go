package cmd

import (
	"fmt"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/logging"
)

type backupSrv interface {
	Backup(sdCardFolder, outputFolder string) (string, error)
}

// backupCmd represents the generate command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup all data and config files of the osmlogger",
	Long:  `backup all data and config files of the open sea map logger as a zip file`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		outputFolder, _ := cmd.Flags().GetString("output")
		return Backup(sdCardFolder, outputFolder)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringP("output", "o", "", "output folder. Default is actual working folder")
}

// Backup get the checker and execute it on the sd file set
func Backup(sdCardFolder, outputFolder string) error {
	bck := do.MustInvokeAs[backupSrv](internal.Inj)
	td := time.Now()
	name, err := bck.Backup(sdCardFolder, outputFolder)
	if JSONOutput {
		fmt.Printf(`{ "filename":"%s" }`, name)
	}
	logging.Root.Infof("backup files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	return err
}
