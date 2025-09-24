package cmd

import (
	"fmt"
	"time"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/check"
	"github.com/willie68/osmltools/internal/logging"
)

// touchCmd touches data files with the right file date
var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "touches the creation time of the data files",
	Long:  `touches the creation time of the data files of the open sea map logger`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		files, _ := cmd.Flags().GetStringSlice("files")
		return Touch(sdCardFolder, files)
	},
}

func init() {
	rootCmd.AddCommand(touchCmd)

	touchCmd.Flags().StringSliceP("files", "f", []string{}, "files to process, separated by commas")
}

// Touch get the checker and execute touch on the sd file set
func Touch(sdCardFolder string, files []string) error {
	chk := do.MustInvoke[check.Checker](internal.Inj)
	td := time.Now()
	res, err := chk.Touch(sdCardFolder, files)
	logging.Root.Infof("touching files took %d seconds", time.Since(td).Abs().Milliseconds()/1000)
	if err == nil {
		if JSONOutput {
			fmt.Println(res.JSON())
			return nil
		}
		if res.Result {
			fmt.Println("touching files successful")
		} else {
			fmt.Println("touching files with errors")
		}
		for _, m := range res.Messages {
			fmt.Print(m)
		}
	}
	return err
}
