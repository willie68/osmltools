package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/export"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var trackCmd = &cobra.Command{
	Use:    "track",
	Short:  "convert the data file(s) to a track and zip it",
	Hidden: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		cmd.Root().SilenceUsage = true
		cmd.Root().SilenceErrors = true
		JSONOutput = true
		logging.Root.SetLevel(logging.None)
		internal.Init()
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		files, _ := cmd.Flags().GetStringSlice("files")
		return Convert(sdCardFolder, files)
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)

	trackCmd.Flags().StringSliceP("files", "f", []string{}, "files to process, separated by commas")
	trackCmd.Flags().StringP("output", "o", "./", "output folder. Default is the working dir.")
}

// Convert get the exporter and execute it on the sd file set
func CreateTrack(sdCardFolder string, files []string) error {
	exp := do.MustInvoke[export.Exporter](nil)
	res, err := exp.Convert(sdCardFolder, files)
	if err != nil {
		return err
	}
	res.LogLines = make([]*model.LogLine, 0)
	js, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Print(string(js))
	return nil
}
