package cmd

import (
	"fmt"
	"runtime"

	"github.com/ndk123-web/trak/internal/config"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

const (
	Build = "2026.08"
)

var versionCmd = cobra.Command{
	Use:   "version",
	Short: "Display the current installed version of Trak",
	Long: `Display detailed version and build information for the Trak CLI,
including Go runtime version, host operating system architecture, and registry configuration.`,
	Example: `  trak version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		fmt.Printf("  %s%s⚡ Trak CLI%s %s(%s)%s\n", ui.Bold, ui.Green, ui.Reset, ui.White, config.TrakConfig.Version, ui.Reset)
		fmt.Printf("  %s──────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)
		fmt.Printf("  %s• Version     :%s  %s\n", ui.Gray, ui.Reset, config.TrakConfig.Version)
		fmt.Printf("  %s• Build       :%s  %s\n", ui.Gray, ui.Reset, Build)
		fmt.Printf("  %s• Go Runtime  :%s  %s\n", ui.Gray, ui.Reset, runtime.Version())
		fmt.Printf("  %s• Platform    :%s  %s/%s\n", ui.Gray, ui.Reset, runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  %s• Registry    :%s  github.com/%s/%s (%s)\n",
			ui.Gray, ui.Reset,
			config.TrakConfig.GithubUsername,
			config.TrakConfig.RegistryName,
			config.TrakConfig.RepositoryBranch,
		)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(&versionCmd)
}
