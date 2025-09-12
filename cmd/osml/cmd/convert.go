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

var convertCmd = &cobra.Command{
	Use:    "convert",
	Short:  "convert the data file(s) to an defined format for the UI",
	Hidden: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Root().SilenceUsage = true
		cmd.Root().SilenceErrors = true
		JSONOutput = true
		logging.Root.SetLevel(logging.None)
		internal.Init()
	},
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
	res, err := exp.Convert(sdCardFolder)
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
