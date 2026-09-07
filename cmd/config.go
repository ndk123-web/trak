package cmd

import "github.com/spf13/cobra"

type User struct {
	UserName string
	Email    string
	Password string
}

var configCmd = cobra.Command{
	Use:     "config",
	Aliases: []string{"set"},
	Short:   "Set Your Username, Password Or Email",
	Args:    cobra.ExactArgs(0),
}

func init() {
	var user User = User{}
	configCmd.Flags().StringVarP(&user.UserName, "username", "u", "guest", "Set the Username For Trak Workspace")
	configCmd.Flags().StringVarP(&user.Email, "email", "e", "guest@gmail.com", "Set the Email")
	configCmd.Flags().StringVarP(&user.Password, "password", "p", "trak@123", "Set the Password")

	rootCmd.AddCommand(&configCmd)
}
