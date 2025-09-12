package cmd

import (
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/export"
	"github.com/willie68/osmltools/internal/logging"
)

var convertCmd = &cobra.Command{
	Use:    "convert",
	Short:  "convert the data file(s) to an defined format for the UI",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Convert(sdCardFolder)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

// Convert get the exporter and execute it on the sd file set
func Convert(sdCardFolder string) error {
	exp := do.MustInvoke[export.Exporter](nil)
	td := time.Now()
	err := exp.Convert(sdCardFolder)
	logging.Root.Infof("converting file took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	return err
}
