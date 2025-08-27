package cmd

import (
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/check"
)

// checkCmd represents the generate command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check the data files of the osmlogger",
	Long:  `check the data files of the open sea map logger`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFile, _ := cmd.Flags().GetString("output")
		return Check(sdCardFolder, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringP("output", "o", "./", "output folder. Default is actual working folder")
}

// Check get the checker and execute it on the sd file set
func Check(sdCardFolder string, outputFolder string) error {
	chk := do.MustInvoke[check.Checker](nil)
	return chk.Check(sdCardFolder, outputFolder)
}
