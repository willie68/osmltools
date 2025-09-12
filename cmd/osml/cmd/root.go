package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
	"github.com/willie68/osmltools/internal/logging"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "osml",
		Short: "open sea map logger tools",
		Long:  `open sea map logger tools`,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			logging.Root.SetLevel(logging.Info)
			if verbose {
				logging.Root.SetLevel(logging.Debug)
			}
			if JSONOutput {
				cmd.Root().SilenceUsage = true
				cmd.Root().SilenceErrors = true
				logging.Root.SetLevel(logging.None)
			}
			internal.Init()
		},
	}
	sdCardFolder string
	verbose      bool
	JSONOutput   bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&sdCardFolder, "sdcard", "s", "./", "root folder of the logger sd card")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&JSONOutput, "json", "", false, "output as json where applicable")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
