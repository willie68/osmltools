package cmd

import (
	"fmt"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/model"
	"github.com/willie68/osmltools/internal/track"
)

var (
	trackCmd = &cobra.Command{
		Use:    "track",
		Short:  "convert the data file(s) to a track and zip it",
		Hidden: false,
	}

	newTrackCmd = &cobra.Command{
		Use:    "new",
		Short:  "create a new track and add data to it",
		Hidden: false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			files, _ := cmd.Flags().GetStringSlice("files")
			trackfile, _ := cmd.Flags().GetString("track")
			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")
			vesselID, _ := cmd.Flags().GetInt32("vesselid")
			t := model.Track{
				Files:       make([]model.SourceData, 0),
				Name:        name,
				Description: description,
				VesselID:    vesselID,
			}
			return NewTrack(sdCardFolder, files, trackfile, t)
		},
	}

	addDataTrackCmd = &cobra.Command{
		Use:    "add",
		Short:  "add data to a track file",
		Hidden: false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			files, _ := cmd.Flags().GetStringSlice("files")
			trackfile, _ := cmd.Flags().GetString("track")
			return AddTrack(sdCardFolder, files, trackfile)
		},
	}

	listTrackCmd = &cobra.Command{
		Use:    "list",
		Short:  "list all information about a track file",
		Hidden: false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			track, _ := cmd.Flags().GetString("track")
			return ListTrack(track)
		},
	}
)

func init() {
	rootCmd.AddCommand(trackCmd)
	trackCmd.PersistentFlags().StringP("track", "t", "", "the track file to work with")

	trackCmd.AddCommand(newTrackCmd)
	newTrackCmd.Flags().StringSliceP("files", "f", []string{}, "files to process, separated by commas")
	newTrackCmd.Flags().StringP("name", "n", "track", "name of the track")
	newTrackCmd.Flags().StringP("description", "d", "", "description of the track")
	newTrackCmd.Flags().Int32P("vesselid", "i", 0, "vessel id")

	trackCmd.AddCommand(addDataTrackCmd)
	addDataTrackCmd.Flags().StringSliceP("files", "f", []string{}, "files to process, separated by commas")

	trackCmd.AddCommand(listTrackCmd)
}

// NewTrack creates a new track file and adds the given data files to it
func NewTrack(sdCardFolder string, files []string, trackfile string, tr model.Track) error {
	tm := do.MustInvokeAs[track.Manager](internal.Inj)
	return tm.NewTrack(sdCardFolder, files, trackfile, tr)
}

// AddTrack add  data files to an existing track file
func AddTrack(sdCardFolder string, files []string, trackfile string) error {
	tm := do.MustInvokeAs[track.Manager](internal.Inj)
	return tm.AddTrack(sdCardFolder, files, trackfile)
}

// ListTrack lists information about the given track file
func ListTrack(trackfile string) error {
	tm := do.MustInvokeAs[track.Manager](internal.Inj)
	tr, err := tm.ListTrack(trackfile)
	if err == nil {
		if JSONOutput {
			js, err := tr.JSON()
			if err != nil {
				return err
			}
			fmt.Println(js)
			return nil
		}
		fmt.Printf("Track: %s\r\n", trackfile)
		fmt.Printf("Name: %s\r\n", tr.Name)
		fmt.Printf("Description: %s\r\n", tr.Description)
		fmt.Printf("VesselID: %d\r\n", tr.VesselID)
		fmt.Printf("Files: \r\n")
		for _, f := range tr.Files {
			fmt.Printf(" - %s (%d) \r\n", f.FileName, f.Size)
		}
	}
	return err
}
