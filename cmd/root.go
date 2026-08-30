package cmd

import (
	"fmt"

	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "trak",
	Short: "Trak - Developer Learning Workspace Generator",
	Long: `  _______ _____            _  __
 |__   __|  __ \     /\   | |/ /
    | |  | |__) |   /  \  | ' / 
    | |  |  _  /   / /\ \ |  <  
    | |  | | \ \  / ____ \| . \ 
    |_|  |_|  \_\/_/    \_\_|\_\

Trak is a developer CLI tool that materializes structured, in-depth learning workspaces
and production-grade curriculum directly onto your filesystem.

Explore 20+ comprehensive curricula across Languages, Operating Systems, Cloud,
Databases, and DevOps tools.`,
	Example: `  # Discover all available curriculum tracks:
  trak list

  # Filter by category:
  trak list lang
  trak list db
  trak list os
  trak list cloud
  trak list tool

  # Initialize a workspace in the current directory:
  trak init lang/go
  trak init db/postgres
  trak init tool/docker --path ./my-docker-lab

  # Check CLI version:
  trak version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n  %s%s⚡ Trak - Developer Learning Workspace Generator%s\n", ui.Bold, ui.Cyan, ui.Reset)
		fmt.Printf("  %sLocal-first curriculum materialization CLI (v%s)%s\n\n", ui.Gray, Version, ui.Reset)

		fmt.Printf("  %sAvailable Commands:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    %slist%s       Explore all available tracks in an interactive tree\n", ui.Green, ui.Reset)
		fmt.Printf("    %sinit%s       Initialize a hands-on learning workspace\n", ui.Green, ui.Reset)
		fmt.Printf("    %sversion%s    Display CLI version and build details\n\n", ui.Green, ui.Reset)

		fmt.Printf("  %sQuick Start:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    %strak list%s                   # Browse all tracks\n", ui.Cyan, ui.Reset)
		fmt.Printf("    %strak init lang/go%s           # Generate Go learning workspace\n", ui.Cyan, ui.Reset)
		fmt.Printf("    %strak init --help%s            # View init options & examples\n\n", ui.Cyan, ui.Reset)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
