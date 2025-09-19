package cmd

import (
	"fmt"
	"os"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/config"
)

var (
	// versionCmd represents the version command

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "the osml version",
		Long:  `the open sea map logger tools version`,
		RunE: func(_ *cobra.Command, _ []string) error {
			ver := do.MustInvoke[config.Version](internal.Inj)
			if JSONOutput {
				v, err := ver.JSON()
				if err != nil {
					return err
				}
				fmt.Println(v)
			} else {
				fmt.Printf("osml: %s \r\n", os.Args[0])
				fmt.Printf("Version: %s, %s, builded %s \r\n", ver.Version(), ver.Commit(), ver.Date())
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
