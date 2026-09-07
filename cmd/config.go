package cmd

import (
	"fmt"

	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var user models.UserConfig = models.UserConfig{}

var configCmd = cobra.Command{
	Use:     "config",
	Aliases: []string{"set"},
	Short:   "Set Your Username, Password Or Email",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ok, err := helper.UpdateUserConfig(&user)
		if err != nil || !ok {
			ui.Error(fmt.Sprintf("Error: %v", err))
			return
		}
	},
}

func init() {
	configCmd.Flags().StringVarP(&user.Username, "username", "u", "", "Set the Username For Trak Workspace")
	configCmd.Flags().StringVarP(&user.Email, "email", "e", "", "Set the Email")
	configCmd.Flags().StringVarP(&user.Password, "password", "p", "", "Set the Password")

	rootCmd.AddCommand(&configCmd)
}
