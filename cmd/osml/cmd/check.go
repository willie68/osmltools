package cmd

import (
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/logging"
)

// checkCmd represents the generate command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check the data files of the osmlogger",
	Long:  `check the data files of the open sea map logger  and write cleanup to an output folder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFile, _ := cmd.Flags().GetString("output")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		return Check(sdCardFolder, outputFile, overwrite)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringP("output", "o", "./", "output folder. Default is actual working folder")
	checkCmd.Flags().BoolP("overwrite", "w", false, "overwrite already converted files. Default false")
}

// Check get the checker and execute it on the sd file set
func Check(sdCardFolder string, outputFolder string, overwrite bool) error {
	chk := do.MustInvoke[check.Checker](nil)
	td := time.Now()
	err := chk.Check(sdCardFolder, outputFolder, overwrite)
	logging.Root.Infof("checking files took %d seconds", time.Since(td))
	return err
}
