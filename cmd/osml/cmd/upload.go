package cmd

import (
	"fmt"

	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/willie68/osmltools/internal"
)

type credManager interface {
	StoreCredentials(user, password string) error
	GetCredentials(user string) (string, error)
}

var (
	// versionCmd represents the version command

	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "upload a track to oseam",
		Long:  `upload a track to oseam using the given credetials`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cred := do.MustInvokeAs[credManager](internal.Inj)
			user, _ := cmd.Flags().GetString("user")
			pwd, _ := cmd.Flags().GetString("password")
			err := cred.StoreCredentials(user, pwd)
			if err != nil {
				return err
			}

			pwd2, err := cred.GetCredentials(user)
			if err != nil {
				return err
			}
			if pwd != pwd2 {
				fmt.Println("Can't use the credentials manager")
				return nil
			}
			fmt.Println("Credentials ok")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP("user", "u", "", "user name of the credentials to use for the upload")
	uploadCmd.Flags().StringP("password", "p", "", "password of the credentials to use for the upload")
}
