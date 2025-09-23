package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/convert"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

var convertCmd = &cobra.Command{
	Use:    "convert",
	Short:  "convert the data file(s) to an defined format for the UI",
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
		track, _ := cmd.Flags().GetString("track")
		return Convert(sdCardFolder, files, track)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringSliceP("files", "f", []string{}, "files to process, separated by commas")
	convertCmd.Flags().StringP("track", "t", "", "the track file to work with")
}

// Convert get the exporter and execute it on the sd file set
func Convert(sdCardFolder string, files []string, track string) error {
	cnv := do.MustInvoke[convert.Converter](internal.Inj)
	res, err := cnv.Convert(sdCardFolder, files, track)
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
