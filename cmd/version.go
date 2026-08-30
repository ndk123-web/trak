package cmd

import (
	"fmt"

	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var versionCmd = cobra.Command{
	Use:   "version",
	Short: "To see the Version Of Tool In Your Machine",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s%s\n", ui.Green, "trak v1.0.0", ui.Reset)
	},
}

func init() {
	rootCmd.AddCommand(&versionCmd)
}
