package cmd

import (
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/export"
	"github.com/willie68/osmltools/internal/logging"
)

// checkCmd represents the generate command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports the data files into one file",
	Long:  `exports the data files of the open sea map logger and write a cleanup version to an output file with the specifig format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFile, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")
		return Export(sdCardFolder, outputFile, format)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP("output", "o", "./track.nmea", "output file. Default is track.nmea in the working folder")
	exportCmd.Flags().StringP("format", "f", "GPX", "the format of the output file. Defaults to GPX, also availble: ")
}

// Export get the exporter and execute it on the sd file set
func Export(sdCardFolder, outputFile, format string) error {
	exp := do.MustInvoke[export.Exporter](nil)
	td := time.Now()
	err := exp.Export(sdCardFolder, outputFile, format)
	logging.Root.Infof("checking files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	return err
}
