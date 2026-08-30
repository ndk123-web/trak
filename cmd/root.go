package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "trak",
	Short: "Learning Tool",
	Long:  "A local-first CLI that creates structured learning workspaces for languages, tools, and technologies.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Trak!")
		fmt.Println("Version v1.0.0")
	},
}

func Execute() error {

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
