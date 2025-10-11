package cmd

import (
	"fmt"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

type checkerSrv interface {
	Check(sdCardFolder, outputFolder string, overwrite, report bool) (*model.CheckResult, error)
}

// checkCmd represents the generate command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check the data files of the osmlogger",
	Long:  `check the data files of the open sea map logger and write report to an output folder`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		outputFolder, _ := cmd.Flags().GetString("output")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		report, _ := cmd.Flags().GetBool("report")
		return Check(sdCardFolder, outputFolder, overwrite, report)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringP("output", "o", "", "output folder. Default is actual working folder")
	checkCmd.Flags().BoolP("overwrite", "w", false, "overwrite already converted files. Default false")
	checkCmd.Flags().BoolP("report", "r", false, "create an report file")
}

// Check get the checker and execute it on the sd file set
func Check(sdCardFolder, outputFolder string, overwrite, report bool) error {
	chk := do.MustInvokeAs[checkerSrv](internal.Inj)
	td := time.Now()
	res, err := chk.Check(sdCardFolder, outputFolder, overwrite, report)
	logging.Root.Infof("checking files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	if err == nil {
		if JSONOutput {
			fmt.Println(res.JSON())
			return nil
		}
	}
	return err
}
