package cmd

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/export"
	"github.com/willie68/osmltools/internal/logging"
)

// checkCmd represents the generate command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports the data files into files",
	Long:  `checks the data files of the open sea map logger, building tracks by day and write a cleanup version to output files with the specifig format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFolder, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")
		format = strings.ToUpper(strings.TrimSpace(format))
		if !slices.Contains(export.SupportedFormats, format) {
			return fmt.Errorf("the format %s is not supported. Supported formats are: %v", format, export.SupportedFormats)
		}
		return Export(sdCardFolder, outputFolder, format)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP("output", "o", "./", "output folder. Default is the working dir. Naming track_####.nmea")
	exportCmd.Flags().StringP("format", "f", export.NMEAFormat, "the format of the output file. Defaults to NMEA, also available: GPX(not implemented yet), KML(not implemented yet)")
}

// Export get the exporter and execute it on the sd file set
func Export(sdCardFolder, outputFolder, format string) error {
	exp := do.MustInvoke[export.Exporter](nil)
	td := time.Now()
	err := exp.Export(sdCardFolder, outputFolder, format)
	logging.Root.Infof("checking files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	return err
}
