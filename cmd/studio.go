package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	targetPort string
)

var studioCmd = cobra.Command{
	Use:   "studio",
	Args:  cobra.ExactArgs(1),
	Short: "Run the Current Workspace in Your Web ",

	Run: func(cmd *cobra.Command, args []string) {
		portToSelect := targetPort
		fmt.Println("Studio: ", portToSelect)
	},
}

func init() {
	studioCmd.Flags().StringVarP(&targetPort, "port", "p", "8200", "Can Give Custom Port On that the Studio Will run")
	rootCmd.AddCommand(&studioCmd)
}
