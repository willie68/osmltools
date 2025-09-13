package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/backup"
	"github.com/willie68/osmltools/internal/logging"
)

// restoreCmd represents the generate command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore all data and config files of the osmlogger",
	Long:  `restore all data and config files of the open sea map logger from a zip file`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		zipfile, _ := cmd.Flags().GetString("zipfile")
		return Restore(zipfile, sdCardFolder)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringP("zipfile", "z", "", "zipfile to restore")
}

// Restore get the checker and execute it on the sd file set
func Restore(zipfile, sdCardFolder string) error {
	bck := do.MustInvoke[backup.Backup](nil)
	td := time.Now()
	zip, err := bck.Restore(zipfile, sdCardFolder)
	if JSONOutput {
		jp := struct {
			Filename string `json:"filename"`
		}{
			Filename: zip,
		}
		js, _ := json.Marshal(jp)
		fmt.Print(string(js))
	}
	logging.Root.Infof("restore files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	return err
}
