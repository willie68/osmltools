package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal/config"
)

var (
	// CmdVersion holding all version information
	CmdVersion config.Version

	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "the osml version",
		Long:  `the open sea map logger tools version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("osml: %s \r\n", os.Args[0])
			fmt.Printf("Version: %s, %s, builded %s \r\n", CmdVersion.Version(), CmdVersion.Commit(), CmdVersion.Date())
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
