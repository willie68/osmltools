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
	RunE: func(cmd *cobra.Command, _ []string) error {
		outputFolder, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")
		name, _ := cmd.Flags().GetString("name")
		format = strings.ToUpper(strings.TrimSpace(format))
		if !slices.Contains(export.SupportedFormats, format) {
			return fmt.Errorf("the format %s is not supported. Supported formats are: %v", format, export.SupportedFormats)
		}
		return Export(sdCardFolder, outputFolder, format, name)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP("output", "o", "./", "output folder. Default is the working dir. Naming track_####.nmea")
	exportCmd.Flags().StringP("format", "f", export.NMEAFormat, "the format of the output file. Defaults to NMEA, also available: GPX, KML, KMZ, GEOJSON")
	exportCmd.Flags().StringP("name", "n", "", "give the track a name")
}

// Export get the exporter and execute it on the sd file set
func Export(sdCardFolder, outputFolder, format, name string) error {
	exp := do.MustInvoke[export.Exporter](nil)
	td := time.Now()
	err := exp.Export(sdCardFolder, outputFolder, format, name)
	logging.Root.Infof("exporting files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	return err
}
