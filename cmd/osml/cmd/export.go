package cmd

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/export"
	"github.com/willie68/osmltools/internal/logging"
	"github.com/willie68/osmltools/internal/model"
)

// checkCmd represents the generate command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports the data files into files",
	Long:  `checks the data files of the open sea map logger, building tracks by day and write a cleanup version to output files with the specifig format`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		output, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")
		files, _ := cmd.Flags().GetStringSlice("files")
		name, _ := cmd.Flags().GetString("name")
		format = strings.ToUpper(strings.TrimSpace(format))
		track, _ := cmd.Flags().GetString("track")

		if !slices.Contains(export.SupportedFormats, format) {
			return fmt.Errorf("the format %s is not supported. Supported formats are: %v", format, export.SupportedFormats)
		}
		if track != "" {
			return ExportTrack(track, output, format)
		}
		return Export(sdCardFolder, output, files, format, name)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringSliceP("files", "f", []string{}, "files to process, separated by commas")
	exportCmd.Flags().StringP("output", "o", "./", "output folder/file. Default is the working dir. Naming track_####.nmea")
	exportCmd.Flags().StringP("format", "m", export.NMEAFormat, "the format of the output file. Defaults to NMEA, also available: GPX, KML, KMZ, GEOJSON")
	exportCmd.Flags().StringP("name", "n", "", "give the track a name")
	exportCmd.Flags().StringP("track", "t", "", "the track file to work with")
}

// Export get the exporter and execute it on the sd file set
func Export(sdCardFolder, outputFolder string, files []string, format, name string) error {
	exp := do.MustInvoke[export.Exporter](internal.Inj)
	td := time.Now()
	err := exp.Export(sdCardFolder, outputFolder, files, format, name)
	logging.Root.Infof("exporting files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	if err == nil {
		if JSONOutput {
			fmt.Println(model.GeneralResult{Result: true}.JSON())
			return nil
		}
		fmt.Println("ok")
	}
	return err
}

// ExportTrack a single track file into the given format
func ExportTrack(trackfile, outputFile, format string) error {
	exp := do.MustInvoke[export.Exporter](internal.Inj)
	td := time.Now()
	err := exp.ExportTrack(trackfile, outputFile, format)
	logging.Root.Infof("exporting track took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	if err == nil {
		if JSONOutput {
			fmt.Println(model.GeneralResult{
				Result:   true,
				Messages: []string{fmt.Sprintf("exported track %s to %s as %s", trackfile, outputFile, format)},
			}.JSON())
			return nil
		}
		fmt.Println("ok")
	}
	return err
}
